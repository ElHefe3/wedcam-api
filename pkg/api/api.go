package api

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"wedcam-api/pkg/nextcloud"

	"github.com/gin-gonic/gin"
)

type AccountResponse struct {
    AccountToken string `json:"account_token"`
    Active      bool   `json:"active"`
}

type QRCodeRequest struct {
    Amount int `json:"amount"`
}

type QRCodeResponse struct {
    QRCodes []QRCode `json:"qr_codes"`
}

type QRCode struct {
    URL         string `json:"url"`
}

type UploadResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

// generateToken creates a random token
func generateToken() (string, error) {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

// CreateAccountHandler generates a new account with a unique token
func CreateAccountHandler(c *gin.Context) {
    token, err := generateToken()
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate account token"})
        return
    }

    // Here you would store the account token in your database
    // with active status = true

    c.JSON(200, AccountResponse{
        AccountToken: token,
        Active:      true,
    })
}

// GenerateQRCodesHandler generates the requested number of QR codes
func GenerateQRCodesHandler(c *gin.Context) {
    accountToken := c.GetHeader("X-Account-Token")
    if accountToken == "" {
        c.JSON(401, gin.H{"error": "Missing account token"})
        return
    }

    // Verify account token and active status here
    // if !isValidAccount(accountToken) {
    //     c.JSON(401, gin.H{"error": "Invalid or inactive account"})
    //     return
    // }

    var req QRCodeRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }

    if req.Amount <= 0 || req.Amount > 100 { // Add reasonable limits
        c.JSON(400, gin.H{"error": "Invalid amount requested"})
        return
    }

    qrCodes := make([]QRCode, req.Amount)
    for i := 0; i < req.Amount; i++ {
        uploadToken, err := generateToken()
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to generate QR codes"})
            return
        }

        // Store the upload token in your database, associated with the account
        
        qrCodes[i] = QRCode{
            URL: os.Getenv("CAMERA_URL")+"?token=" + uploadToken,
        }
    }

    c.JSON(200, QRCodeResponse{
        QRCodes: qrCodes,
    })
}

func ImageUploadHandler(c *gin.Context) {
    uploadToken := c.GetHeader("X-Upload-Token")
    if uploadToken == "" {
        c.JSON(401, UploadResponse{
            Success: false,
            Message: "Missing upload token",
        })
        return
    }

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

    // Store the image with association to the upload token
    if err := nextcloud.UploadImage(header.Filename, buffer); err != nil {
        c.JSON(500, UploadResponse{
            Success: false,
            Message: err.Error(),
        })
        return
    }

    // Mark the upload token as used in your database

    c.JSON(200, UploadResponse{
        Success: true,
        Message: "File uploaded successfully",
    })
}