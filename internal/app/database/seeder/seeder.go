package seeder

import (
	"time"

	"github.com/vizucode/concurent-domain-checker/internal/app/dto/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) (err error) {

	data := &models.DomainCheckHistory{
		Name:      "test",
		Total:     5,
		Success:   5,
		Failed:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Domains: []models.Domain{
			{
				BatchId:     1,
				FullUrl:     "https://google.com",
				StatusCode:  200,
				RedirectUrl: "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				BatchId:     1,
				FullUrl:     "https://facebook.com",
				StatusCode:  200,
				RedirectUrl: "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				BatchId:     1,
				FullUrl:     "https://twitter.com",
				StatusCode:  200,
				RedirectUrl: "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				BatchId:     1,
				FullUrl:     "https://linkedin.com",
				StatusCode:  200,
				RedirectUrl: "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				BatchId:     1,
				FullUrl:     "https://github.com",
				StatusCode:  200,
				RedirectUrl: "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}

	if err := db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}
