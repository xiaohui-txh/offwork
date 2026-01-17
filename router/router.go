package router

import (
	"offwork-backend/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/offwork/checkin", handler.OffworkCheckin)
		v1.GET("/offwork/nearby", handler.NearbyOffwork2)
	}

	return r
}
