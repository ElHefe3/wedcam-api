package server

import (
	"wedcam-api/pkg/nextcloud"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Define routes
	r.GET("/nextcloud/files", nextcloud.GetFilesHandler)

	return r
}

func Run() error {
	// Initialize the router
	router := setupRouter()

	// Start the server
	return router.Run(":8080")
}
