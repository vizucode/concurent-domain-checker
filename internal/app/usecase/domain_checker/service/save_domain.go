package service

import (
	"context"

	"github.com/vizucode/concurent-domain-checker/internal/app/dto/models"
)

func (s *domainCheckerService) saveDomain(ctx context.Context, domain *models.DomainCheckHistory) error {
	err := s.database.CreateDomainHistory(ctx, domain)
	if err != nil {
		return err
	}

	return nil
}
