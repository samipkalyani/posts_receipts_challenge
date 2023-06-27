package router

import (
	handler "github.com/samipkalyani/posts_receipts_challenge/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	handler.SetupPointsMap()

	router.GET("/receipts/:id/points", handler.GetPoints)
	router.POST("/receipts/process", handler.PostReceipts)

	return router
}
