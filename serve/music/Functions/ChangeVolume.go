package Functions

import (
	"MusicBot/config"
	"MusicBot/serve/music"
	"fmt"
)

func ChangeVolume(volume int) {
	logger := config.Logger
	if <-music.PlayStatus.FFmpegSignel == music.STOP {
		music.PlayStatus.Volume = volume
		music.PlayStatus.FFmpegSignel <- music.STOP
	} else {
		music.StopFFmpeg(music.PlayStatus.FFmpegCmd)
		music.PlayStatus.Volume = volume
		<-music.PlayStatus.FFmpegSignel
		music.RunFFmpeg(music.PlayStatus.Music.File)
	}
	logger.Info().Msg(fmt.Sprintf(">>> Change volume to %d <<<", volume))
}
