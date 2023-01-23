package QQ

import (
	"MusicBot/config"
	"MusicBot/serve/music"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type MusicPlayURL struct {
	URL   string `json:"url"`
	Error bool   `json:"error"`
}

type MusicResp struct {
	Data struct {
		PlayURL map[string]MusicPlayURL `json:"playUrl"`
	} `json:"data"`
}

type MusicReq struct {
	ID string
}

type MusicInfoResp struct {
	Response struct {
		Code     int    `json:"code"`
		Ts       int64  `json:"ts"`
		StartTs  int64  `json:"start_ts"`
		Traceid  string `json:"traceid"`
		Songinfo struct {
			Code int `json:"code"`
			Data struct {
				Info struct {
					Company struct {
						Title   string `json:"title"`
						Type    string `json:"type"`
						Content []struct {
							ID        int    `json:"id"`
							Value     string `json:"value"`
							Mid       string `json:"mid"`
							Type      int    `json:"type"`
							ShowType  int    `json:"show_type"`
							IsParent  int    `json:"is_parent"`
							Picurl    string `json:"picurl"`
							ReadCnt   int    `json:"read_cnt"`
							Author    string `json:"author"`
							Jumpurl   string `json:"jumpurl"`
							OriPicurl string `json:"ori_picurl"`
						} `json:"content"`
						Pos         int    `json:"pos"`
						More        int    `json:"more"`
						Selected    string `json:"selected"`
						UsePlatform int    `json:"use_platform"`
					} `json:"company"`
					Genre struct {
						Title   string `json:"title"`
						Type    string `json:"type"`
						Content []struct {
							ID        int    `json:"id"`
							Value     string `json:"value"`
							Mid       string `json:"mid"`
							Type      int    `json:"type"`
							ShowType  int    `json:"show_type"`
							IsParent  int    `json:"is_parent"`
							Picurl    string `json:"picurl"`
							ReadCnt   int    `json:"read_cnt"`
							Author    string `json:"author"`
							Jumpurl   string `json:"jumpurl"`
							OriPicurl string `json:"ori_picurl"`
						} `json:"content"`
						Pos         int    `json:"pos"`
						More        int    `json:"more"`
						Selected    string `json:"selected"`
						UsePlatform int    `json:"use_platform"`
					} `json:"genre"`
					Lan struct {
						Title   string `json:"title"`
						Type    string `json:"type"`
						Content []struct {
							ID        int    `json:"id"`
							Value     string `json:"value"`
							Mid       string `json:"mid"`
							Type      int    `json:"type"`
							ShowType  int    `json:"show_type"`
							IsParent  int    `json:"is_parent"`
							Picurl    string `json:"picurl"`
							ReadCnt   int    `json:"read_cnt"`
							Author    string `json:"author"`
							Jumpurl   string `json:"jumpurl"`
							OriPicurl string `json:"ori_picurl"`
						} `json:"content"`
						Pos         int    `json:"pos"`
						More        int    `json:"more"`
						Selected    string `json:"selected"`
						UsePlatform int    `json:"use_platform"`
					} `json:"lan"`
					PubTime struct {
						Title   string `json:"title"`
						Type    string `json:"type"`
						Content []struct {
							ID        int    `json:"id"`
							Value     string `json:"value"`
							Mid       string `json:"mid"`
							Type      int    `json:"type"`
							ShowType  int    `json:"show_type"`
							IsParent  int    `json:"is_parent"`
							Picurl    string `json:"picurl"`
							ReadCnt   int    `json:"read_cnt"`
							Author    string `json:"author"`
							Jumpurl   string `json:"jumpurl"`
							OriPicurl string `json:"ori_picurl"`
						} `json:"content"`
						Pos         int    `json:"pos"`
						More        int    `json:"more"`
						Selected    string `json:"selected"`
						UsePlatform int    `json:"use_platform"`
					} `json:"pub_time"`
				} `json:"info"`
				Extras struct {
					Name      string `json:"name"`
					Transname string `json:"transname"`
					Subtitle  string `json:"subtitle"`
					From      string `json:"from"`
					Wikiurl   string `json:"wikiurl"`
				} `json:"extras"`
				TrackInfo struct {
					ID       int    `json:"id"`
					Type     int    `json:"type"`
					Mid      string `json:"mid"`
					Name     string `json:"name"`
					Title    string `json:"title"`
					Subtitle string `json:"subtitle"`
					Singer   []struct {
						ID    int    `json:"id"`
						Mid   string `json:"mid"`
						Name  string `json:"name"`
						Title string `json:"title"`
						Type  int    `json:"type"`
						Uin   int    `json:"uin"`
					} `json:"singer"`
					Album struct {
						ID         int    `json:"id"`
						Mid        string `json:"mid"`
						Name       string `json:"name"`
						Title      string `json:"title"`
						Subtitle   string `json:"subtitle"`
						TimePublic string `json:"time_public"`
						Pmid       string `json:"pmid"`
					} `json:"album"`
					Mv struct {
						ID    int    `json:"id"`
						Vid   string `json:"vid"`
						Name  string `json:"name"`
						Title string `json:"title"`
						Vt    int    `json:"vt"`
					} `json:"mv"`
					Interval   int    `json:"interval"`
					Isonly     int    `json:"isonly"`
					Language   int    `json:"language"`
					Genre      int    `json:"genre"`
					IndexCd    int    `json:"index_cd"`
					IndexAlbum int    `json:"index_album"`
					TimePublic string `json:"time_public"`
					Status     int    `json:"status"`
					Fnote      int    `json:"fnote"`
					File       struct {
						MediaMid      string        `json:"media_mid"`
						Size24Aac     int           `json:"size_24aac"`
						Size48Aac     int           `json:"size_48aac"`
						Size96Aac     int           `json:"size_96aac"`
						Size192Ogg    int           `json:"size_192ogg"`
						Size192Aac    int           `json:"size_192aac"`
						Size128Mp3    int           `json:"size_128mp3"`
						Size320Mp3    int           `json:"size_320mp3"`
						SizeApe       int           `json:"size_ape"`
						SizeFlac      int           `json:"size_flac"`
						SizeDts       int           `json:"size_dts"`
						SizeTry       int           `json:"size_try"`
						TryBegin      int           `json:"try_begin"`
						TryEnd        int           `json:"try_end"`
						URL           string        `json:"url"`
						SizeHires     int           `json:"size_hires"`
						HiresSample   int           `json:"hires_sample"`
						HiresBitdepth int           `json:"hires_bitdepth"`
						B30S          int           `json:"b_30s"`
						E30S          int           `json:"e_30s"`
						Size96Ogg     int           `json:"size_96ogg"`
						Size360Ra     []interface{} `json:"size_360ra"`
						SizeDolby     int           `json:"size_dolby"`
						SizeNew       []int         `json:"size_new"`
					} `json:"file"`
					Pay struct {
						PayMonth   int `json:"pay_month"`
						PriceTrack int `json:"price_track"`
						PriceAlbum int `json:"price_album"`
						PayPlay    int `json:"pay_play"`
						PayDown    int `json:"pay_down"`
						PayStatus  int `json:"pay_status"`
						TimeFree   int `json:"time_free"`
					} `json:"pay"`
					Action struct {
						Switch   int `json:"switch"`
						Msgid    int `json:"msgid"`
						Alert    int `json:"alert"`
						Icons    int `json:"icons"`
						Msgshare int `json:"msgshare"`
						Msgfav   int `json:"msgfav"`
						Msgdown  int `json:"msgdown"`
						Msgpay   int `json:"msgpay"`
						Switch2  int `json:"switch2"`
						Icon2    int `json:"icon2"`
					} `json:"action"`
					Ksong struct {
						ID  int    `json:"id"`
						Mid string `json:"mid"`
					} `json:"ksong"`
					Volume struct {
						Gain float64 `json:"gain"`
						Peak int     `json:"peak"`
						Lra  float64 `json:"lra"`
					} `json:"volume"`
					Label       string   `json:"label"`
					URL         string   `json:"url"`
					Bpm         int      `json:"bpm"`
					Version     int      `json:"version"`
					Trace       string   `json:"trace"`
					DataType    int      `json:"data_type"`
					ModifyStamp int      `json:"modify_stamp"`
					Pingpong    string   `json:"pingpong"`
					Ppurl       string   `json:"ppurl"`
					Tid         int      `json:"tid"`
					Ov          int      `json:"ov"`
					Sa          int      `json:"sa"`
					Es          string   `json:"es"`
					Vs          []string `json:"vs"`
					Vi          []int    `json:"vi"`
				} `json:"track_info"`
			} `json:"data"`
		} `json:"songinfo"`
	} `json:"response"`
}

func QueryMusic(id string) (*music.Music, error) {
	logger := config.Logger
	if music.Musics.GetMusicByID(id) != nil {
		logger.Info().Msg(fmt.Sprintf("Music %s added from cache", music.Musics.GetMusicByID(id).Name))
		return music.Musics.GetMusicByID(id), nil
	}
	url, err := QueryMusicURL(id)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to query music url")
		return nil, err
	}
	path, err := DownloadMusic(id, url)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to download music")
		return nil, err
	}
	music, err := QueryMusicInfo(id)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to query music info")
		return nil, err
	}
	music.File = path
	logger.Info().Msgf("Music %s downloaded", music.Name)
	return music, nil
}

func QueryMusicURL(id string) (string, error) {
	logger := config.Logger
	musicReq := &MusicReq{
		ID: id,
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.Config.QQAPI+"/getMusicPlay", nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create request")
		return "", err
	}
	//req.Header.Set("Cookie", config.Config.NetEaseCookie)
	params := req.URL.Query()
	params.Add("songmid", musicReq.ID)
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to send request")
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read response")
		return "", err
	}
	var musicResp MusicResp
	logger.Debug().Msg(string(body))
	if err := json.Unmarshal(body, &musicResp); err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal response")
		return "", err
	}
	return musicResp.Data.PlayURL[musicReq.ID].URL, nil
}

func DownloadMusic(id string, url string) (string, error) {
	logger := config.Logger
	resp, err := http.Get(url)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to download music")
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error().Msg("Failed to download music, status code: " + strconv.Itoa(resp.StatusCode))
		return "", err
	}
	defer resp.Body.Close()
	f, err := os.Create("./assets/music/Q" + id + ".mp3")
	defer f.Close()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create file")
		return "", err
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to write file")
		return "", err
	}
	return "./assets/music/Q" + id + ".mp3", nil
}

func QueryMusicInfo(id string) (*music.Music, error) {
	logger := config.Logger
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.Config.QQAPI+"/getSongInfo", nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create request")
		return nil, err
	}
	//req.Header.Set("Cookie", config.Config.NetEaseCookie)
	params := req.URL.Query()
	params.Add("songmid", id)
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to send request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read response")
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logger.Error().Msg("Failed to get music info")
		return nil, err
	}
	var musicInfoResp MusicInfoResp
	if err := json.Unmarshal(body, &musicInfoResp); err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal response")
		return nil, err
	}
	if musicInfoResp.Response.Code != 0 {
		logger.Error().Msg("Failed to get music info")
		return nil, err
	}

	ar := make([]string, 0)
	for _, v := range musicInfoResp.Response.Songinfo.Data.TrackInfo.Singer {
		ar = append(ar, v.Name)
	}
	return &music.Music{
		ID:       id,
		Name:     musicInfoResp.Response.Songinfo.Data.TrackInfo.Name,
		Artists:  ar,
		Album:    "https://i2.hdslb.com/bfs/face/29acac2dd587c7dd4ca85f93b4d080fb17cfb401.jpg",
		LastTime: musicInfoResp.Response.Songinfo.Data.TrackInfo.Interval * 1000,
	}, nil
}
