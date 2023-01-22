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
	logger.Info().Msgf(">>> Play music: %s from %s to %s <<<", music.Name, music.File, MusicStatus.DstAddr)
	MusicStatus.Music = music
	s := <-MusicStatus.KOOKSignel
	if s == STOP {
		RunKOOKVoice(MusicStatus.Channel)
	} else {
		StopKOOKVoice(MusicStatus.KOOKCmd)
		<-MusicStatus.KOOKSignel
		RunKOOKVoice(MusicStatus.Channel)
	}
	s = <-MusicStatus.FFmpegSignel
	if s == STOP {
		RunFFmpeg(music.File)
	} else {
		StopFFmpeg(MusicStatus.FFmpegCmd)
		<-MusicStatus.FFmpegSignel
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
		MusicStatus.DstAddr,
		"-t",
		config.Config.BotToken,
	)
	var err error
	go func() {
		MusicStatus.KOOKSignel <- RUN
		MusicStatus.KOOKCmd = cmd
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
	MusicStatus.KOOKSignel <- STOP
}

func ChangeChannel(channel int64) {
	logger := config.Logger
	if <-MusicStatus.KOOKSignel == STOP {
		MusicStatus.Channel = channel
		MusicStatus.KOOKSignel <- STOP
	} else {
		StopKOOKVoice(MusicStatus.KOOKCmd)
		MusicStatus.Channel = channel
		<-MusicStatus.KOOKSignel
		RunKOOKVoice(MusicStatus.Channel)
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
		fmt.Sprintf("volume=%ddB", MusicStatus.Volume),
		"-ab",
		"120k",
		"-acodec",
		"libopus",
		"-f",
		"mpegts",
		MusicStatus.DstAddr,
	)
	go func() {
		MusicStatus.FFmpegSignel <- RUN
		MusicStatus.FFmpegCmd = cmd
		//out, err := cmd.CombinedOutput()
		err := cmd.Run()
		if err != nil {
			logger.Error().Err(err).Msg(fmt.Sprintf("Failed to run ffmpeg, pid: %d", cmd.Process.Pid))
			//logger.Error().Msg(string(out))
		}
		logger.Info().Msg(fmt.Sprintf(">>> FFmpeg finished, pid: %d <<<", cmd.Process.Pid))
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
	MusicStatus.FFmpegSignel <- STOP
}

func ChangeVolume(volume int) {
	logger := config.Logger
	if <-MusicStatus.FFmpegSignel == STOP {
		MusicStatus.Volume = volume
		MusicStatus.FFmpegSignel <- STOP
	} else {
		StopFFmpeg(MusicStatus.FFmpegCmd)
		MusicStatus.Volume = volume
		<-MusicStatus.FFmpegSignel
		RunFFmpeg(MusicStatus.Music.File)
	}
	logger.Info().Msg(fmt.Sprintf(">>> Change volume to %d <<<", volume))
}
