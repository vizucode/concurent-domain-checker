package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vizucode/concurent-domain-checker/configs/database"
	"github.com/vizucode/concurent-domain-checker/internal/app/routes"
	"github.com/vizucode/concurent-domain-checker/internal/app/usecase/domain_checker/controllers"
	"github.com/vizucode/concurent-domain-checker/internal/app/usecase/domain_checker/repository"
	"github.com/vizucode/concurent-domain-checker/internal/app/usecase/domain_checker/service"
)

func Run() {
	router := gin.Default()

	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := database.NewDatabaseConnection(dsn)
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	httpClient := http.Client{
		Timeout: 2 * time.Minute,
		Transport: &http.Transport{
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 150,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	// Dependency Injection
	databaseRepo := repository.NewDatabaseRepository(db)
	domainService := service.NewDomainCheckerService(databaseRepo, &httpClient)
	domainController := controllers.NewDomainCheckerController(domainService)

	// Register Routes
	routes.NewRoute(router, domainController)

	router.Run(":8080")
}
