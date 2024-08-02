package music

import (
	"MusicBot/config"
	"fmt"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// ffmpeg -re -i ./res/${music.id}.mp3 -af volume=${this.status.volume}dB -ab 120k -acodec libopus -f mpegts zmq:tcp://127.0.0.1:5559
func PlayMusic(musicReq *Music) {
	logger := config.Logger
	logger.Info().Msgf(">>> Play musicReq: %s from %s to %s <<<", musicReq.Name, musicReq.File, PlayStatus.DstAddr)
	PlayStatus.Music = musicReq
	s := <-PlayStatus.KOOKSignel
	if s == STOP {
		RunKOOKVoice(PlayStatus.Channel)
	} else {
		StopKOOKVoice(PlayStatus.KOOKCmd)
		<-PlayStatus.KOOKSignel
		RunKOOKVoice(PlayStatus.Channel)
	}

	s = <-PlayStatus.FFmpegSignel
	if musicReq == nil {
		return
	}
	if s == STOP {
		RunFFmpeg(musicReq.File)
	} else {
		StopFFmpeg(PlayStatus.FFmpegCmd)
		<-PlayStatus.FFmpegSignel
		RunFFmpeg(musicReq.File)
	}
	logger.Info().Msg(fmt.Sprintf(">>> Start playing %s <<<", musicReq.Name))
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
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

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
	time.Sleep(3 * time.Second)
	err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	logger.Info().Msg(fmt.Sprintf(">>> Stop KOOKVoice, pid: %d <<<", cmd.Process.Pid))
	if err != nil {
		logger.Error().Err(err).Msg(fmt.Sprintf("Failed to stop KOOKVoice, pid: %d", cmd.Process.Pid))
	}
	PlayStatus.KOOKSignel <- STOP
}

func RunFFmpeg(src string) *exec.Cmd {
	logger := config.Logger
	time.Sleep(3 * time.Second)
	logger.Info().Msg(">>> Run FFmpeg <<<")
	cmd := exec.Command("ffmpeg",
		"-re",
		"-i",
		src,
		"-af",
		fmt.Sprintf("loudnorm=i=-16,volume=%ddB", PlayStatus.Volume),
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
		out, err := cmd.CombinedOutput()
		// err := cmd.Run()
		if err != nil {
			logger.Error().Err(err).Msg(fmt.Sprintf("Failed to run ffmpeg, pid: %d", cmd.Process.Pid))
			logger.Error().Msg(cmd.String())
			logger.Error().Msg(string(out))
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
