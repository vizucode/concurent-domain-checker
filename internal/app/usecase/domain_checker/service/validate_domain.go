package service

import (
	"context"
	"net/http"
	"sync"

	"github.com/vizucode/concurent-domain-checker/internal/app/dto/models"
)

func (s *domainCheckerService) checkDomain(ctx context.Context, domains <-chan string) <-chan models.Domain {
	var (
		result    = make(chan models.Domain)
		wg        = sync.WaitGroup{}
		maxWorker = 150
	)

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
						result <- models.Domain{
							FullUrl: url,
						}
						continue
					}

					resp, err := s.apiClient.Do(req)
					if err != nil {
						result <- models.Domain{
							FullUrl: url,
						}
						continue
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
