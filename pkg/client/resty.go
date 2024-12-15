// pkg/client/resty.go
package clients

import (
    "github.com/go-resty/resty/v2"
)

var Client *resty.Client

func InitResty() {
    Client = resty.New().
        SetBaseURL("https://your-nextcloud-instance.com").
        SetHeader("Content-Type", "application/json").
        SetTimeout(10 * 1000 * 1000)
}