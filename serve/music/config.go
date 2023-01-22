package music

import (
	"MusicBot/config"
	"fmt"
	"os"
	"os/exec"
)

type Music struct {
	ID       int
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
	Channel       int64
	DstAddr       string
	Volume        int
	KOOKRunning   bool
	KOOKCmd       *exec.Cmd
	KOOKSignel    chan bool
	FFmpegRunning bool
	FFmpegCmd     *exec.Cmd
	FFmpegSignel  chan bool
	Music         *Music
}

const (
	STOP = false
	RUN  = true
)

var Musics MusicsList
var MusicStatus Status

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
	MusicStatus.Channel = config.Config.VoiceChannel[0].ID
	MusicStatus.DstAddr = fmt.Sprintf("zmq:tcp://127.0.0.1:%d", config.Config.VoicePort)
	MusicStatus.Volume = -20
	MusicStatus.KOOKRunning = false
	MusicStatus.FFmpegRunning = false
	MusicStatus.KOOKSignel = make(chan bool, 1)
	MusicStatus.FFmpegSignel = make(chan bool, 1)
	MusicStatus.KOOKSignel <- STOP
	MusicStatus.FFmpegSignel <- STOP
	return nil
}
