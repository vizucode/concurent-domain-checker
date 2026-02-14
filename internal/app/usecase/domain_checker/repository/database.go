package repository

import (
	"context"

	"github.com/vizucode/concurent-domain-checker/internal/app/dto/models"
	"gorm.io/gorm"
)

type DatabaseRepository interface {
	CreateDomainHistory(ctx context.Context, history *models.DomainCheckHistory) error
}

type databaseRepository struct {
	db *gorm.DB
}

func NewDatabaseRepository(db *gorm.DB) DatabaseRepository {
	return &databaseRepository{
		db: db,
	}
}

func (r *databaseRepository) CreateDomainHistory(ctx context.Context, history *models.DomainCheckHistory) error {
	if err := r.db.WithContext(ctx).Create(history).Error; err != nil {
		return err
	}

	return nil
}
