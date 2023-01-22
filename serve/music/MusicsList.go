package music

func (m *MusicsList) Add(music *Music) {
	m.Musics = append(m.Musics, *music)
}

func (m *MusicsList) GetMusicByID(id int) *Music {
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
