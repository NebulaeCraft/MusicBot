package music

import "os"

type Music struct {
	ID      int
	Name    string
	Artists []string
	Album   string
	File    string
}

type MusicsList struct {
	Musics []Music
}

func InitMusicEnv() error {
	err := os.RemoveAll("./assets/music")
	if err != nil {
		return err
	}
	err = os.Mkdir("./assets/music", 0755)
	if err != nil {
		return err
	}
	return nil
}
