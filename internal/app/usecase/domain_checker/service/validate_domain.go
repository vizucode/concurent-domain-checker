package service

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vizucode/concurent-domain-checker/internal/app/dto/models"
)

type Metrics struct {
	TotalProcessed atomic.Int64
	TotalSuccess   atomic.Int64
	TotalFailed    atomic.Int64
	ProcessingTime time.Duration
}

func (s *domainCheckerService) checkDomain(ctx context.Context, domains <-chan string, metrics *Metrics) <-chan models.Domain {
	var (
		result    = make(chan models.Domain)
		wg        = sync.WaitGroup{}
		maxWorker = 150
	)

	slog.Info("Starting worker pool", "workers", maxWorker)

	for range maxWorker {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case url, ok := <-domains:
					if !ok {
						return
					}

					req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
					if err != nil {
						slog.Error("Failed to check domain", "url", url, "error", err)
						metrics.TotalFailed.Add(1)
						metrics.TotalProcessed.Add(1)
						result <- models.Domain{
							FullUrl:    url,
							StatusCode: 503,
						}
						continue
					}

					resp, err := s.apiClient.Do(req)
					if err != nil {
						slog.Error("Failed to check domain", "url", url, "error", err)
						metrics.TotalFailed.Add(1)
						metrics.TotalProcessed.Add(1)
						result <- models.Domain{
							FullUrl:    url,
							StatusCode: 503,
						}
						continue
					}

					metrics.TotalProcessed.Add(1)
					if resp.StatusCode >= 200 && resp.StatusCode < 400 {
						metrics.TotalSuccess.Add(1)
					} else {
						metrics.TotalFailed.Add(1)
					}

					redirectURL := ""
					if resp.Request != nil && resp.Request.URL.String() != url {
						redirectURL = resp.Request.URL.String()
					}

					resp.Body.Close()

					result <- models.Domain{
						FullUrl:     url,
						StatusCode:  resp.StatusCode,
						RedirectUrl: redirectURL,
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}
