package Bili

import (
	"MusicBot/config"
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

func QueryBiliAudio(serial string, isBV bool) (*music.Music, error) {
	logger := config.Logger
	if isBV {
		DownloadVideoAudio("BV" + serial)
	} else {
		DownloadVideoAudio("AV" + serial)
	}
	videoInfo, err := QueryVideoInfo(serial, isBV)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to query video info")
		return nil, err
	}
	path, err := ChangeFileName(serial, videoInfo.Title)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to change filename")
		return nil, err
	}

	musicResp := &music.Music{
		ID:       serial,
		Name:     videoInfo.Title,
		Artists:  []string{videoInfo.Up},
		Album:    videoInfo.Cover,
		File:     path,
		LastTime: videoInfo.Duration * 1000,
	}

	return musicResp, nil

}
