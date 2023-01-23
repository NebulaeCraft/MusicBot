package music

import (
	"github.com/lonelyevil/kook"
)

func (m *MusicsList) Add(M *Music) {
	m.Musics = append(m.Musics, *M)
}

func (m *MusicsList) GetMusicByID(id string) *Music {
	for _, music := range m.Musics {
		if music.ID == id {
			return &music
		}
	}
	return nil
}

func (m *MusicsList) GetMusicByName(name string) *Music {
	for _, music := range m.Musics {
		if music.Name == name {
			return &music
		}
	}
	return nil
}

func (m *MusicsList) Play(ctx *kook.KmarkdownMessageContext) {
	s := <-PlayStatus.PlaySignel
	if s == STOP {
		go PlayMusic(&m.Musics[0])
		SendMusicCard(ctx, &m.Musics[0])
		m.Musics = m.Musics[1:]
	}
}
