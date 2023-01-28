package bff

import (
	"MusicBot/serve/music/Functions"
	"github.com/gin-gonic/gin"
	"github.com/lonelyevil/kook"
	"net/http"
)

func SkipHandler(ctx *kook.KmarkdownMessageContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		Functions.SkipMusic()
		c.JSON(http.StatusOK, gin.H{"status": 200, "message": "Skip successfully!"})
	}
}
