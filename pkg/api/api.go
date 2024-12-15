// pkg/api/api.go
package api

import (
    "github.com/gin-gonic/gin"
    "wedcam-api/pkg/nextcloud"
)

type UploadResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

func ImageUploadHandler(c *gin.Context) {
    // Get file from request
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(400, UploadResponse{
            Success: false,
            Message: "Failed to read image file",
        })
        return
    }

    // Open the file
    src, err := file.Open()
    if err != nil {
        c.JSON(500, UploadResponse{
            Success: false,
            Message: "Failed to open image file",
        })
        return
    }
    defer src.Close()

    // Read file into memory
    buffer := make([]byte, file.Size)
    if _, err := src.Read(buffer); err != nil {
        c.JSON(500, UploadResponse{
            Success: false,
            Message: "Failed to read image data",
        })
        return
    }

    // Upload to Nextcloud
    if err := nextcloud.UploadImage(file.Filename, buffer); err != nil {
        c.JSON(500, UploadResponse{
            Success: false,
            Message: err.Error(),
        })
        return
    }

    c.JSON(200, UploadResponse{
        Success: true,
        Message: "Image uploaded successfully",
    })
}