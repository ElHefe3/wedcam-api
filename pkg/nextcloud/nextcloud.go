// pkg/nextcloud/nextcloud.go
package nextcloud

import (
    "fmt"
    "wedcam-api/pkg/client"
)

func UploadImage(fileName string, imageData []byte) error {
    resp, err := clients.Client.R().
        SetBody(imageData).
        SetHeader("Content-Type", "image/jpeg").
        Put("/remote.php/dav/files/guest/" + fileName)
    
    if err != nil {
        return fmt.Errorf("failed to upload image: %v", err)
    }

    if resp.StatusCode() != 201 && resp.StatusCode() != 204 {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
    }

    return nil
}
