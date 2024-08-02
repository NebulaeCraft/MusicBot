package music

import (
	"github.com/lonelyevil/kook"
	"time"
)

func (m *MusicsList) Add(M *Music) {
	if PlayStatus.CanAppend == false {
		SendMsg(PlayStatus.Ctx, "当前播放列表已锁定，无法添加新的音乐")
		return
	}
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

func (m *MusicsList) GetIndexByID(id string) int {
	for i, music := range m.Musics {
		if music.ID == id {
			return i
		}
	}
	return -1
}

func (m *MusicsList) Play(ctx *kook.KmarkdownMessageContext) {
	s := <-PlayStatus.PlaySignel
	if s == STOP {
		time.Sleep(5 * time.Second)
		go PlayMusic(&m.Musics[0])
		SendMusicCard(ctx, &m.Musics[0])
		m.Musics = m.Musics[1:]
	}
}

func (m *MusicsList) PlayBtn(ctx *kook.MessageButtonClickContext) {
	s := <-PlayStatus.PlaySignel
	if s == STOP {
		time.Sleep(5 * time.Second)
		go PlayMusic(&m.Musics[0])
		SendMusicCard(PlayStatus.Ctx, &m.Musics[0])
		m.Musics = m.Musics[1:]
	}
}
