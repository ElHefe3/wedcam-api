package nextcloud

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFilesHandler(c *gin.Context) {
	// Call Nextcloud service
	data, err := FetchFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{"data": data})
}
