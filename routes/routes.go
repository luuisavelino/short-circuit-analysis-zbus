package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luuisavelino/short-circuit-analysis-zbus/controllers"
	"github.com/luuisavelino/short-circuit-analysis-zbus/middleware"
)

func HandleRequest() {
	router := gin.New()

	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/actuator/health"),
		gin.Recovery(),
		middleware.Logger(),
	)

	actuator := router.Group("/actuator/health")
	{
		actuator.GET("/", controllers.HealthGET)
	}

	zbus := router.Group("/api/v2/files/:fileId/zbus")
	{
		zbus.GET("/", controllers.AllZbus)
		zbus.GET("/:seq", controllers.ZbusSeq)
		zbus.GET("/atuacao/:line", controllers.AllZbus)
		zbus.GET("/atuacao/:line/bars", controllers.Bars)
		zbus.GET("/atuacao/:line/bars/:bar", controllers.Bar)
		zbus.GET("/short-circuit/:line/point/:point", controllers.AllZbus)
		zbus.GET("/short-circuit/:line/point/:point/bars", controllers.Bars)
		zbus.GET("/short-circuit/:line/point/:point/bars/:bar", controllers.Bar)
	}

	router.Run(":8081")
}
