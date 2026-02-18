package service

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/vizucode/concurent-domain-checker/internal/app/dto/domains"
	"github.com/vizucode/concurent-domain-checker/internal/app/dto/models"
	"github.com/vizucode/concurent-domain-checker/internal/app/usecase/domain_checker/repository"
)

type DomainCheckerService interface {
	RequestDomain(ctx context.Context, request *domains.DomainCheckerRequest) (domains.DomainCheckerResponse, error)
}

type domainCheckerService struct {
	database  repository.DatabaseRepository
	apiClient *http.Client
	jobCh     chan models.JobRequest
}

func NewDomainCheckerService(database repository.DatabaseRepository, apiClient *http.Client) DomainCheckerService {
	svc := &domainCheckerService{
		database:  database,
		apiClient: apiClient,
		jobCh:     make(chan models.JobRequest, 150),
	}

	maxWorker := 150

	go func() {
		for range maxWorker {
			go svc.worker()
		}
	}()

	return svc
}

func (s *domainCheckerService) worker() {
	for job := range s.jobCh {
		func() {
			defer job.Wg.Done()

			req, err := http.NewRequestWithContext(job.Ctx, http.MethodGet, job.Url, nil)
			if err != nil {
				slog.Error("Failed to check domain", "url", job.Url, "error", err)
				job.Metrics.TotalFailed.Add(1)
				job.Metrics.TotalProcessed.Add(1)
				job.Result <- models.Domain{
					FullUrl:    job.Url,
					StatusCode: 503,
				}
				return
			}

			resp, err := s.apiClient.Do(req)
			if err != nil {
				slog.Error("Failed to check domain", "url", job.Url, "error", err)
				job.Metrics.TotalFailed.Add(1)
				job.Metrics.TotalProcessed.Add(1)
				job.Result <- models.Domain{
					FullUrl:    job.Url,
					StatusCode: 503,
				}
				return
			}

			job.Metrics.TotalProcessed.Add(1)
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				job.Metrics.TotalSuccess.Add(1)
			} else {
				job.Metrics.TotalFailed.Add(1)
			}

			redirectURL := ""
			if resp.Request != nil && resp.Request.URL.String() != job.Url {
				redirectURL = resp.Request.URL.String()
			}

			resp.Body.Close()

			job.Result <- models.Domain{
				FullUrl:     job.Url,
				StatusCode:  resp.StatusCode,
				RedirectUrl: redirectURL,
			}
		}()
	}
}

func (s *domainCheckerService) RequestDomain(ctx context.Context, request *domains.DomainCheckerRequest) (domains.DomainCheckerResponse, error) {

	var (
		domainCheckHistory = models.DomainCheckHistory{}
		metrics            = &models.Metrics{}
	)

	start := time.Now()
	chanDomain := s.sanitizeDomain(request.Domains)

	chanCheckDomain := s.checkDomain(ctx, chanDomain, metrics)

	domainCheckHistory.Name = request.Name

	for domain := range chanCheckDomain {

		domainCheckHistory.Domains = append(domainCheckHistory.Domains, models.Domain{
			FullUrl:     domain.FullUrl,
			StatusCode:  domain.StatusCode,
			RedirectUrl: domain.RedirectUrl,
		})

		domainCheckHistory.Total++

		if domain.StatusCode >= 200 && domain.StatusCode < 400 {
			domainCheckHistory.Success++
		} else {
			domainCheckHistory.Failed++
		}
	}

	metrics.ProcessingTime = time.Since(start)

	err := s.saveDomain(ctx, &domainCheckHistory)
	if err != nil {
		return domains.DomainCheckerResponse{
			Message: "Failed to save domain",
			Success: false,
		}, nil
	}

	return domains.DomainCheckerResponse{
		Message: "Success",
		Success: true,
	}, nil
}
