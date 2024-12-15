package api

import (
	"fmt"
	"net/http"
)

func ImageUploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Image upload endpoint - Not implemented yet!")
}
