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
		fmt.Sprintf("https://www.bilibili.com/video/%s", serial),
		"--only-audio",
		"--dir",
		"./assets/music",
	)
	fmt.Println(cmd.String())
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
		//if strings.HasPrefix(file.Name(), title[:int(len(title)/2)]) {
		if !(strings.HasPrefix(file.Name(), "Q") || strings.HasPrefix(file.Name(), "B") || strings.HasPrefix(file.Name(), "N") || strings.HasPrefix(file.Name(), "U")) {
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

	serial = strings.Split(serial, "/")[0]
	videoInfo, err := QueryVideoInfo(strings.Split(serial, "?")[0], isBV)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to query video info")
		return nil, err
	}
	path, err := ChangeFileName(strings.Replace(serial, "/", ".", -1), videoInfo.Title)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to change filename, filename: " + videoInfo.Title)
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
