package bff

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func UploadHandler(dst string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _ := c.FormFile("file")
		if strings.HasSuffix(file.Filename, ".mp3") || strings.HasSuffix(file.Filename, ".wav") || strings.HasSuffix(file.Filename, ".ogg") {
			err := c.SaveUploadedFile(file, dst+"/"+file.Filename)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to save file")
			}
			logger.Info().Msg(fmt.Sprintf("File %s uploaded successfully", file.Filename))
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "File uploaded successfully!"})
		} else {
			logger.Info().Msg(fmt.Sprintf("File %s 's type is not supported", file.Filename))
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "File type not supported!"})
		}
	}
}
