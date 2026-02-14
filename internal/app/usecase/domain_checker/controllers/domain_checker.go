package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vizucode/concurent-domain-checker/internal/app/dto/domains"
	"github.com/vizucode/concurent-domain-checker/internal/app/usecase/domain_checker/service"
)

type DomainCheckerController interface {
	RequestDomain(c *gin.Context)
}

type domainCheckerController struct {
	service service.DomainCheckerService
}

func NewDomainCheckerController(service service.DomainCheckerService) DomainCheckerController {
	return &domainCheckerController{
		service: service,
	}
}

func (ctrl *domainCheckerController) RequestDomain(c *gin.Context) {
	var request domains.DomainCheckerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.service.RequestDomain(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
