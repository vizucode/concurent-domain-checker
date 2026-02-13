package seeder

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/vizucode/concurent-domain-checker/internal/app/dto/models"
)

func Seed(db *gorm.DB) (err error) {

	data := &models.DomainCheckHistory{
		ID:        1,
		Name:      "test",
		Total:     5,
		Success:   5,
		Failed:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Domains: []models.Domain{
			{
				Id:          1,
				BatchId:     1,
				FullUrl:     "https://google.com",
				StatusCode:  "200",
				RedirectUrl: "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Id:          2,
				BatchId:     1,
				FullUrl:     "https://facebook.com",
				StatusCode:  "200",
				RedirectUrl: "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Id:          3,
				BatchId:     1,
				FullUrl:     "https://twitter.com",
				StatusCode:  "200",
				RedirectUrl: "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Id:          4,
				BatchId:     1,
				FullUrl:     "https://linkedin.com",
				StatusCode:  "200",
				RedirectUrl: "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Id:          5,
				BatchId:     1,
				FullUrl:     "https://github.com",
				StatusCode:  "200",
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
