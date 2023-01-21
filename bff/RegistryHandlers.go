package bff

import (
	"MusicBot/config"
	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
)

var logger *log.Logger

func TestHandler(s *gin.Engine) {
	s.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})
}
func IUpload(s *gin.Engine, dst string) {
	s.POST("/upload", UploadHandler(dst))
}

func RegistryHandlers(s *gin.Engine) {
	logger = config.Logger
	logger.Info().Msg("Registering Gin handlers")
	s.MaxMultipartMemory = 16 << 20 // 8 MiB
	TestHandler(s)
	IUpload(s, "./assets/music")
}
