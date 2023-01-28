package Functions

import (
	"MusicBot/config"
	"MusicBot/serve/music"
	"fmt"
)

func ChangeChannel(channel int64) {
	logger := config.Logger
	if <-music.PlayStatus.KOOKSignel == music.STOP {
		music.PlayStatus.Channel = channel
		music.PlayStatus.KOOKSignel <- music.STOP
	} else {
		music.StopKOOKVoice(music.PlayStatus.KOOKCmd)
		music.PlayStatus.Channel = channel
		<-music.PlayStatus.KOOKSignel
		_, err := music.RunKOOKVoice(music.PlayStatus.Channel)
		if err != nil {
			return
		}
	}
	logger.Info().Msg(fmt.Sprintf(">>> Change channel to %d <<<", channel))
}
