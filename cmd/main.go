package main

import (
    "github.com/gin-gonic/gin"
    "wedcam-api/pkg/client"
    "wedcam-api/pkg/api"
)

func main() {
    clients.InitResty()

    r := gin.Default()

    r.Static("/static", "./static")

    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })

    // Routes
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "Service is running",
        })
    })

    r.POST("/upload", api.ImageUploadHandler)

    r.Run(":8900")
}