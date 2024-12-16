package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	
	"wedcam-api/pkg/db"
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
    URL string `json:"url"`
}

type UploadResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

func generateToken() (string, error) {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func CreateAccountHandler(c *gin.Context) {
    token, err := generateToken()
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate account token"})
        return
    }

    result, err := db.DB.Exec(
        "INSERT INTO accounts (token, active) VALUES (?, ?)",
        token, true)
    if err != nil {
		fmt.Println(err.Error())
        c.JSON(500, gin.H{"error": "Failed to create account"})
        return
    }

    _, err = result.LastInsertId()
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to get account ID"})
        return
    }

    c.JSON(200, AccountResponse{
        AccountToken: token,
        Active:      false,
    })
}

func GenerateQRCodesHandler(c *gin.Context) {
    accountToken := c.GetHeader("X-Account-Token")
    if accountToken == "" {
        c.JSON(401, gin.H{"error": "Missing account token"})
        return
    }

    var accountID int64
    var active bool
    err := db.DB.QueryRow(
        "SELECT id, active FROM accounts WHERE token = ?",
        accountToken).Scan(&accountID, &active)

    if err == sql.ErrNoRows {
        c.JSON(401, gin.H{"error": "Invalid account"})
        return
    } else if err != nil {
        c.JSON(500, gin.H{"error": "Database error"})
        return
    }

    if !active {
        c.JSON(401, gin.H{"error": "Account is inactive"})
        return
    }

    var req QRCodeRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }

    if req.Amount <= 0 || req.Amount > 100 {
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

        _, err = db.DB.Exec(
            "INSERT INTO qr_codes (token, account_id, uploads_allowed, uploads_used) VALUES (?, ?, ?, ?)",
            uploadToken, accountID, 1, 0)
        if err != nil {
            c.JSON(500, gin.H{"error": "Failed to store QR code"})
            return
        }

        qrCodes[i] = QRCode{
            URL: os.Getenv("CAMERA_URL") + "?token=" + uploadToken,
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

    var uploadsUsed int
    var uploadsAllowed int
    err := db.DB.QueryRow(
        "SELECT uploads_used, uploads_allowed FROM qr_codes WHERE token = ?",
        uploadToken).Scan(&uploadsUsed, &uploadsAllowed)

    if err == sql.ErrNoRows {
        c.JSON(401, UploadResponse{
            Success: false,
            Message: "Invalid upload token",
        })
        return
    } else if err != nil {
        c.JSON(500, UploadResponse{
            Success: false,
            Message: "Database error",
        })
        return
    }

    if uploadsUsed >= uploadsAllowed {
        c.JSON(400, UploadResponse{
            Success: false,
            Message: "Upload token has already been used",
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

    if err := nextcloud.UploadImage(header.Filename, buffer); err != nil {
        c.JSON(500, UploadResponse{
            Success: false,
            Message: err.Error(),
        })
        return
    }

    _, err = db.DB.Exec(
        "UPDATE qr_codes SET uploads_used = uploads_used + 1 WHERE token = ?",
        uploadToken)
    if err != nil {
        c.JSON(500, UploadResponse{
            Success: false,
            Message: "Failed to update upload status",
        })
        return
    }

    c.JSON(200, UploadResponse{
        Success: true,
        Message: "File uploaded successfully",
    })
}