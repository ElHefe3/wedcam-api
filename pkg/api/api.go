package api

import (
	"wedcam-api/pkg/nextcloud"
	"github.com/gin-gonic/gin"
)

type UploadResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

func ImageUploadHandler(c *gin.Context) {
    file, header, err := c.Request.FormFile("image")
    if err != nil {
        c.JSON(400, UploadResponse{
            Success: false,
            Message: "Failed to get file from form",
        })
        return
    }
    defer file.Close()

    buffer := make([]byte, header.Size)
    if _, err := file.Read(buffer); err != nil {
        c.JSON(500, UploadResponse{
            Success: false,
            Message: "Failed to read image data",
        })
        return
    }

    if err := nextcloud.UploadImage(header.Filename, buffer); err != nil {
        c.JSON(500, UploadResponse{
            Success: false,
            Message: err.Error(),
        })
        return
    }

    c.JSON(200, UploadResponse{
        Success: true,
        Message: "File uploaded successfully",
    })
}