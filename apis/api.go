package apis

import (
	"log"

	"github.com/gin-gonic/gin"

	"tourist-api/utils"
)

var frontendDomain string

func InitAPI() {
	router := gin.Default()

	router.Use(CORSMiddleware())

	apiList(router)

	err := router.Run(":" + utils.ENV_API_PORT)
	if err != nil {
		log.Fatal(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", utils.ENV_FRONTEND_DOMAIN)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func apiList(router *gin.Engine) {
	router.POST("/api/v1/login", login)

	router.GET("/api/v1/trips", getTrips)
	router.POST("/api/v1/trips", createTrip)
}
