package clients

import (
	"os"
	"time"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)


var Client *resty.Client
func InitResty() {
    godotenv.Load()
    Client = resty.New().
        SetBaseURL(os.Getenv("NC_PUBLIC_URL")).
        SetTimeout(30 * time.Second)
}