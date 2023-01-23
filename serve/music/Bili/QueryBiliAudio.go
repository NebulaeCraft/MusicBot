package Bili

import (
	"MusicBot/serve/music"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func DownloadVideoAudio(serial string) {
	cmd := exec.Command("bilix",
		"get_video",
		"https://www.bilibili.com/video/"+serial,
		"--only-audio",
		"--dir",
		"./assets/music",
	)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func ChangeFileName(serial, title string) (string, error) {
	dir, err := ioutil.ReadDir("./assets/music")
	if err != nil {
		return "", err
	}
	for _, file := range dir {
		if strings.HasPrefix(file.Name(), title[:int(len(title)/2)]) {
			os.Rename("./assets/music/"+file.Name(), "./assets/music/B"+serial+".mp3")
			return "./assets/music/B" + serial + ".mp3", nil
		}
	}
	return "", errors.New("No such file")
}

func QueryBiliAudio(serial string) (*music.Music, error) {
	DownloadVideoAudio(serial)
	filename, err := QueryVideoTitle(serial)
	if err != nil {
		fmt.Println(err)
	}
	path, err := ChangeFileName(serial, filename)
	if err != nil {
		fmt.Println(err)
	}

	musicResp := &music.Music{
		ID:       serial,
		Name:     filename,
		Artists:  nil,
		Album:    "https://i2.hdslb.com/bfs/face/29acac2dd587c7dd4ca85f93b4d080fb17cfb401.jpg",
		File:     path,
		LastTime: 0,
	}

	return musicResp, nil

}
