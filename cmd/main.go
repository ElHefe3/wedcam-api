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
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Account-Token, X-Upload-Token")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })

    // Health check route
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "Service is running",
        })
    })

    // Account and QR code routes
    r.POST("/accounts", api.CreateAccountHandler)
    r.POST("/qr-codes", api.GenerateQRCodesHandler)
    
    // Image upload route
    r.POST("/upload", api.ImageUploadHandler)

    // Start the server
    r.Run(":8900")
}