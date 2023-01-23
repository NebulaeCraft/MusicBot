package music

import (
	"MusicBot/config"
	"fmt"
	"os/exec"
	"strconv"
)

// ffmpeg -re -i ./res/${music.id}.mp3 -af volume=${this.status.volume}dB -ab 120k -acodec libopus -f mpegts zmq:tcp://127.0.0.1:5559
func PlayMusic(music *Music) {
	logger := config.Logger
	logger.Info().Msgf(">>> Play music: %s from %s to %s <<<", music.Name, music.File, PlayStatus.DstAddr)
	PlayStatus.Music = music
	s := <-PlayStatus.KOOKSignel
	if s == STOP {
		RunKOOKVoice(PlayStatus.Channel)
	} else {
		StopKOOKVoice(PlayStatus.KOOKCmd)
		<-PlayStatus.KOOKSignel
		RunKOOKVoice(PlayStatus.Channel)
	}
	s = <-PlayStatus.FFmpegSignel
	if s == STOP {
		RunFFmpeg(music.File)
	} else {
		StopFFmpeg(PlayStatus.FFmpegCmd)
		<-PlayStatus.FFmpegSignel
		RunFFmpeg(music.File)
	}
	logger.Info().Msg(fmt.Sprintf(">>> Start playing %s <<<", music.Name))
}

func RunKOOKVoice(channel int64) (*exec.Cmd, error) {
	logger := config.Logger
	logger.Info().Msg(">>> Run KOOK Voice <<<")
	cmd := exec.Command(config.Config.KOOKVoice,
		"-c",
		strconv.Itoa(int(channel)),
		"-i",
		PlayStatus.DstAddr,
		"-t",
		config.Config.BotToken,
	)
	var err error
	go func() {
		PlayStatus.KOOKSignel <- RUN
		PlayStatus.KOOKCmd = cmd
		err = cmd.Run()
		if err != nil {
			logger.Error().Err(err).Msg(fmt.Sprintf("Failed to run KOOKVoice, pid: %d", cmd.Process.Pid))
		}
		logger.Info().Msg(fmt.Sprintf(">>> KOOKVoice finished, pid: %d <<<", cmd.Process.Pid))
	}()
	return cmd, err
}

func StopKOOKVoice(cmd *exec.Cmd) {
	logger := config.Logger
	err := cmd.Process.Kill()
	logger.Info().Msg(fmt.Sprintf(">>> Stop KOOKVoice, pid: %d <<<", cmd.Process.Pid))
	if err != nil {
		logger.Error().Err(err).Msg(fmt.Sprintf("Failed to stop KOOKVoice, pid: %d", cmd.Process.Pid))
	}
	PlayStatus.KOOKSignel <- STOP
}

func ChangeChannel(channel int64) {
	logger := config.Logger
	if <-PlayStatus.KOOKSignel == STOP {
		PlayStatus.Channel = channel
		PlayStatus.KOOKSignel <- STOP
	} else {
		StopKOOKVoice(PlayStatus.KOOKCmd)
		PlayStatus.Channel = channel
		<-PlayStatus.KOOKSignel
		RunKOOKVoice(PlayStatus.Channel)
	}
	logger.Info().Msg(fmt.Sprintf(">>> Change channel to %d <<<", channel))
}

func RunFFmpeg(src string) *exec.Cmd {
	logger := config.Logger
	logger.Info().Msg(">>> Run FFmpeg <<<")
	cmd := exec.Command("ffmpeg",
		"-re",
		"-i",
		src,
		"-af",
		fmt.Sprintf("volume=%ddB", PlayStatus.Volume),
		"-ab",
		"120k",
		"-acodec",
		"libopus",
		"-f",
		"mpegts",
		PlayStatus.DstAddr,
	)
	go func() {
		PlayStatus.FFmpegSignel <- RUN
		PlayStatus.FFmpegCmd = cmd
		//out, err := cmd.CombinedOutput()
		err := cmd.Run()
		if err != nil {
			logger.Error().Err(err).Msg(fmt.Sprintf("Failed to run ffmpeg, pid: %d", cmd.Process.Pid))
			//logger.Error().Msg(string(out))
		} else {
			logger.Info().Msg(fmt.Sprintf(">>> FFmpeg finished, pid: %d <<<", cmd.Process.Pid))
			PlayStatus.PlaySignel <- STOP
		}
	}()
	return cmd
}

func StopFFmpeg(cmd *exec.Cmd) {
	logger := config.Logger
	err := cmd.Process.Kill()
	logger.Info().Msg(fmt.Sprintf(">>> Stop FFmpeg, pid: %d <<<", cmd.Process.Pid))
	if err != nil {
		logger.Error().Err(err).Msg(fmt.Sprintf("Failed to stop FFmpeg, pid: %d", cmd.Process.Pid))
	}
	PlayStatus.FFmpegSignel <- STOP
}

func ChangeVolume(volume int) {
	logger := config.Logger
	if <-PlayStatus.FFmpegSignel == STOP {
		PlayStatus.Volume = volume
		PlayStatus.FFmpegSignel <- STOP
	} else {
		StopFFmpeg(PlayStatus.FFmpegCmd)
		PlayStatus.Volume = volume
		<-PlayStatus.FFmpegSignel
		RunFFmpeg(PlayStatus.Music.File)
	}
	logger.Info().Msg(fmt.Sprintf(">>> Change volume to %d <<<", volume))
}
