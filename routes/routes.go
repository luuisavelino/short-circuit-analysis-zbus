package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luuisavelino/short-circuit-analysis-zbus/controllers"
	"github.com/luuisavelino/short-circuit-analysis-zbus/middleware"
)

func HandleRequest() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/actuator/health"),
		gin.Recovery(),
		middleware.Logger(),
	)

	actuator := router.Group("/actuator")
	{
		actuator.GET("/health", controllers.HealthGET)
	}

	zbus := router.Group("/api/v2/files/:fileId")
	{
		zbus.GET("/zbus", controllers.AllZbus)
		zbus.GET("/zbus/:seq", controllers.ZbusSeq)
	}

	router.Run(":8081")
}
