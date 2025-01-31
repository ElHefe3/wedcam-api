package nextcloud

import (
    "bytes"
    "fmt"
    "net/http"
    "path"
    "strings"
    "time"
    "wedcam-api/pkg/client"
)

func UploadImage(token string, data []byte) error {
    filename := time.Now().Format("20060102_150405") + ".jpg"

    // Ensure the token is URL-safe
    safeToken := strings.Trim(token, "/")

    // Construct the folder path
    folderPath := "/" + safeToken + "/"

    // Create the folder using MKCOL
    folderResp, err := clients.Client.R().
        SetHeader("Content-Type", "application/xml").
        SetBody("").
        Execute("MKCOL", folderPath) // Use "MKCOL" as a string

    if err != nil {
        return fmt.Errorf("failed to create folder: %w", err)
    }

    // Check if the folder was created successfully
    if folderResp.StatusCode() != http.StatusCreated && folderResp.StatusCode() != http.StatusMethodNotAllowed {
        return fmt.Errorf("failed to create folder, status code %d: %s",
            folderResp.StatusCode(), folderResp.String())
    }

    // Construct the full file path
    filePath := path.Join(folderPath, filename)

    // Upload the file
    fileResp, err := clients.Client.R().
        SetBody(bytes.NewReader(data)).
        Put(filePath)

    if err != nil {
        return fmt.Errorf("failed to upload image: %w", err)
    }

    if fileResp.StatusCode() >= 400 {
        return fmt.Errorf("upload failed with status code %d: %s",
            fileResp.StatusCode(), fileResp.String())
    }

    return nil
}

