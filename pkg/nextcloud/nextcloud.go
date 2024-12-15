package nextcloud

import (
    "bytes"
    "fmt"
    "wedcam-api/pkg/client"
)

func UploadImage(filename string, data []byte) error {
    resp, err := clients.Client.R().
        SetBody(bytes.NewReader(data)).
        Put("/" + filename)

    if err != nil {
        return fmt.Errorf("failed to upload image: %w", err)
    }

    if resp.StatusCode() >= 400 {
        return fmt.Errorf("upload failed with status code %d: %s", 
            resp.StatusCode(), resp.String())
    }

    return nil
}