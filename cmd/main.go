package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "wedcam-api/pkg/client"
    "wedcam-api/pkg/api"
    "wedcam-api/pkg/db"
)

func main() {
    if err := db.InitDB(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.DB.Close()

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

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "Service is running",
        })
    })

    r.POST("/accounts", api.CreateAccountHandler)
    r.POST("/qr-codes", api.GenerateQRCodesHandler)

    r.POST("/upload", api.ImageUploadHandler)

    r.Run(":8900")
}