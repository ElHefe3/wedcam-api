package nextcloud

import (
	"wecam-api/pkg/clients"
)

func FetchFiles() (string, error) {
	resp, err := clients.Client.R().
		Get("/remote.php/dav/files/username/")

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}
