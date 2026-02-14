package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vizucode/concurent-domain-checker/internal/app/usecase/domain_checker/controllers"
)

func NewRoute(route *gin.Engine, controller controllers.DomainCheckerController) {
	route.POST("/domain-checker", controller.RequestDomain)
}
