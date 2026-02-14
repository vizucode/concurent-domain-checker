package service

import (
	"context"
	"net/http"

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
}

func NewDomainCheckerService(database repository.DatabaseRepository, apiClient *http.Client) DomainCheckerService {
	return &domainCheckerService{
		database:  database,
		apiClient: apiClient,
	}
}

func (s *domainCheckerService) RequestDomain(ctx context.Context, request *domains.DomainCheckerRequest) (domains.DomainCheckerResponse, error) {

	var (
		domainCheckHistory = models.DomainCheckHistory{}
	)

	chanDomain := s.sanitizeDomain(request.Domains)

	chanCheckDomain := s.checkDomain(ctx, chanDomain)

	domainCheckHistory.Name = request.Name

	for domain := range chanCheckDomain {

		domainCheckHistory.Domains = append(domainCheckHistory.Domains, models.Domain{
			FullUrl:     domain.FullUrl,
			StatusCode:  domain.StatusCode,
			RedirectUrl: domain.RedirectUrl,
		})

		domainCheckHistory.Total++

		if domain.StatusCode < 400 {
			domainCheckHistory.Success++
		} else {
			domainCheckHistory.Failed++
		}
	}

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
