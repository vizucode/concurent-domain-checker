package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/vizucode/concurent-domain-checker/internal/app/dto/models"
)

func (s *domainCheckerService) checkDomain(ctx context.Context, domains <-chan string, metrics *models.Metrics) <-chan models.Domain {
	var (
		maxWorker = 150
		result    = make(chan models.Domain, 150)
		wg        = &sync.WaitGroup{}
	)

	slog.Info("Starting worker pool", "workers", maxWorker)

	go func() {
		for {
			select {
			case <-ctx.Done():
				goto done
			case url, ok := <-domains:
				if !ok {
					goto done
				}

				wg.Add(1)
				s.jobCh <- models.JobRequest{
					Ctx:     ctx,
					Url:     url,
					Metrics: metrics,
					Result:  result,
					Wg:      wg,
				}
			}
		}

	done:
		wg.Wait()
		close(result)
	}()

	return result
}
