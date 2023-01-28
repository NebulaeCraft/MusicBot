package Functions

import (
	"MusicBot/config"
	"MusicBot/serve/music"
)

func SkipMusic() {
	logger := config.Logger
	if music.PlayStatus.Music == nil {
		logger.Info().Msg("当前无音乐播放")
		return
	}
	<-music.PlayStatus.KOOKSignel
	music.StopKOOKVoice(music.PlayStatus.KOOKCmd)
	<-music.PlayStatus.FFmpegSignel
	music.StopFFmpeg(music.PlayStatus.FFmpegCmd)
	music.PlayStatus.PlaySignel <- music.STOP
	if len(music.Musics.Musics) == 0 {
		return
	}
}
