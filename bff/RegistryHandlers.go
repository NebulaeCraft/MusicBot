package bff

import (
	"MusicBot/config"
	"MusicBot/serve/music"
	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
	"net/http"
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

func ISkip(s *gin.Engine) {
	s.GET("/skip", SkipHandler(music.PlayStatus.Ctx))
}

func IListMusic(s *gin.Engine, Musics *music.MusicsList) {
	s.GET("/music/list", func(c *gin.Context) {
		c.JSON(200, Musics)
	})
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func RegistryHandlers(s *gin.Engine) {
	logger = config.Logger
	logger.Info().Msg("Registering Gin handlers")
	s.Use(Cors())
	s.MaxMultipartMemory = 16 << 20 // 8 MiB
	TestHandler(s)
	IUpload(s, "./assets/music")
	IListMusic(s, &music.Musics)
	ISkip(s)
}
