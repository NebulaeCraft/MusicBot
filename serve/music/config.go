package music

import (
	"MusicBot/config"
	"fmt"
	"os"
	"os/exec"
)

type Music struct {
	ID       string
	Name     string
	Artists  []string
	Album    string
	File     string
	LastTime int
}

type MusicsList struct {
	Musics []Music
}

type Status struct {
	Channel      int64
	DstAddr      string
	Volume       int
	KOOKCmd      *exec.Cmd
	KOOKSignel   chan bool
	FFmpegCmd    *exec.Cmd
	FFmpegSignel chan bool
	PlaySignel   chan bool
	Music        *Music
}

const (
	STOP = false
	RUN  = true
)

var Musics MusicsList
var PlayStatus Status

func InitMusicEnv() error {
	Musics.Musics = make([]Music, 0)
	err := os.RemoveAll("./assets/music")
	if err != nil {
		return err
	}
	err = os.Mkdir("./assets/music", 0755)
	if err != nil {
		return err
	}
	PlayStatus.Channel = config.Config.VoiceChannel[0].ID
	PlayStatus.DstAddr = fmt.Sprintf("zmq:tcp://127.0.0.1:%d", config.Config.VoicePort)
	PlayStatus.Volume = -20
	PlayStatus.KOOKSignel = make(chan bool, 1)
	PlayStatus.FFmpegSignel = make(chan bool, 1)
	PlayStatus.PlaySignel = make(chan bool, 1)
	PlayStatus.KOOKSignel <- STOP
	PlayStatus.FFmpegSignel <- STOP
	PlayStatus.PlaySignel <- STOP
	return nil
}
