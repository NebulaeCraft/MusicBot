package bff

import (
	"MusicBot/serve/music"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func UploadHandler(dst string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _ := c.FormFile("file")
		if strings.HasSuffix(file.Filename, ".mp3") || strings.HasSuffix(file.Filename, ".wav") || strings.HasSuffix(file.Filename, ".ogg") {
			fileName := "U" + time.Now().Format("20060102150405") + ".mp3"
			err := c.SaveUploadedFile(file, dst+"/"+fileName)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to save file")
			}
			logger.Info().Msg(fmt.Sprintf("File %s uploaded successfully", file.Filename))
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "File uploaded successfully!", "filename": fileName})
			musicResult := &music.Music{
				ID:       fileName,
				Name:     file.Filename,
				Artists:  []string{"群星"},
				Album:    "https://i2.hdslb.com/bfs/face/29acac2dd587c7dd4ca85f93b4d080fb17cfb401.jpg",
				File:     dst + "/" + fileName,
				LastTime: 114514,
			}
			music.SendMsg(music.PlayStatus.Ctx, fmt.Sprintf("%s 已加入播放列表", musicResult.Name))
			music.Musics.Add(musicResult)
			go music.Musics.Play(music.PlayStatus.Ctx)
		} else {
			logger.Info().Msg(fmt.Sprintf("File %s 's type is not supported", file.Filename))
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "File type not supported!"})
		}
	}
}
