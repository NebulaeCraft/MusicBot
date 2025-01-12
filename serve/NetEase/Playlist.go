package NetEase

import (
	"MusicBot/config"
	"MusicBot/serve/player"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type PlaylistResp struct {
	Songs []struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
		Pst  int    `json:"pst"`
		T    int    `json:"t"`
		Ar   []struct {
			ID    int           `json:"id"`
			Name  string        `json:"name"`
			Tns   []interface{} `json:"tns"`
			Alias []interface{} `json:"alias"`
		} `json:"ar"`
		Alia []string    `json:"alia"`
		Pop  int         `json:"pop"`
		St   int         `json:"st"`
		Rt   interface{} `json:"rt"`
		Fee  int         `json:"fee"`
		V    int         `json:"v"`
		Crbt interface{} `json:"crbt"`
		Cf   string      `json:"cf"`
		Al   struct {
			ID     int           `json:"id"`
			Name   string        `json:"name"`
			PicURL string        `json:"picUrl"`
			Tns    []interface{} `json:"tns"`
			PicStr string        `json:"pic_str"`
			Pic    int64         `json:"pic"`
		} `json:"al,omitempty"`
		Dt int `json:"dt"`
		H  struct {
			Br   int `json:"br"`
			Fid  int `json:"fid"`
			Size int `json:"size"`
			Vd   int `json:"vd"`
			Sr   int `json:"sr"`
		} `json:"h"`
		M struct {
			Br   int `json:"br"`
			Fid  int `json:"fid"`
			Size int `json:"size"`
			Vd   int `json:"vd"`
			Sr   int `json:"sr"`
		} `json:"m"`
		L struct {
			Br   int `json:"br"`
			Fid  int `json:"fid"`
			Size int `json:"size"`
			Vd   int `json:"vd"`
			Sr   int `json:"sr"`
		} `json:"l"`
		Sq struct {
			Br   int `json:"br"`
			Fid  int `json:"fid"`
			Size int `json:"size"`
			Vd   int `json:"vd"`
			Sr   int `json:"sr"`
		} `json:"sq"`
		Hr                   interface{}   `json:"hr"`
		A                    interface{}   `json:"a"`
		Cd                   string        `json:"cd"`
		No                   int           `json:"no"`
		RtURL                interface{}   `json:"rtUrl"`
		Ftype                int           `json:"ftype"`
		RtUrls               []interface{} `json:"rtUrls"`
		DjID                 int           `json:"djId"`
		Copyright            int           `json:"copyright"`
		SID                  int           `json:"s_id"`
		Mark                 int64         `json:"mark"`
		OriginCoverType      int           `json:"originCoverType"`
		OriginSongSimpleData interface{}   `json:"originSongSimpleData"`
		TagPicList           interface{}   `json:"tagPicList"`
		ResourceState        bool          `json:"resourceState"`
		Version              int           `json:"version"`
		SongJumpInfo         interface{}   `json:"songJumpInfo"`
		EntertainmentTags    interface{}   `json:"entertainmentTags"`
		AwardTags            interface{}   `json:"awardTags"`
		Single               int           `json:"single"`
		NoCopyrightRcmd      interface{}   `json:"noCopyrightRcmd"`
		Mv                   int           `json:"mv"`
		Rtype                int           `json:"rtype"`
		Rurl                 interface{}   `json:"rurl"`
		Mst                  int           `json:"mst"`
		Cp                   int           `json:"cp"`
		PublishTime          int64         `json:"publishTime"`
		Tns                  []string      `json:"tns,omitempty"`
	} `json:"songs"`
	Privileges []struct {
		ID                 int         `json:"id"`
		Fee                int         `json:"fee"`
		Payed              int         `json:"payed"`
		St                 int         `json:"st"`
		Pl                 int         `json:"pl"`
		Dl                 int         `json:"dl"`
		Sp                 int         `json:"sp"`
		Cp                 int         `json:"cp"`
		Subp               int         `json:"subp"`
		Cs                 bool        `json:"cs"`
		Maxbr              int         `json:"maxbr"`
		Fl                 int         `json:"fl"`
		Toast              bool        `json:"toast"`
		Flag               int         `json:"flag"`
		PreSell            bool        `json:"preSell"`
		PlayMaxbr          int         `json:"playMaxbr"`
		DownloadMaxbr      int         `json:"downloadMaxbr"`
		MaxBrLevel         string      `json:"maxBrLevel"`
		PlayMaxBrLevel     string      `json:"playMaxBrLevel"`
		DownloadMaxBrLevel string      `json:"downloadMaxBrLevel"`
		PlLevel            string      `json:"plLevel"`
		DlLevel            string      `json:"dlLevel"`
		FlLevel            string      `json:"flLevel"`
		Rscl               interface{} `json:"rscl"`
		FreeTrialPrivilege struct {
			ResConsumable      bool        `json:"resConsumable"`
			UserConsumable     bool        `json:"userConsumable"`
			ListenType         int         `json:"listenType"`
			CannotListenReason int         `json:"cannotListenReason"`
			PlayReason         interface{} `json:"playReason"`
			FreeLimitTagType   interface{} `json:"freeLimitTagType"`
		} `json:"freeTrialPrivilege"`
		RightSource    int `json:"rightSource"`
		ChargeInfoList []struct {
			Rate          int         `json:"rate"`
			ChargeURL     interface{} `json:"chargeUrl"`
			ChargeMessage interface{} `json:"chargeMessage"`
			ChargeType    int         `json:"chargeType"`
		} `json:"chargeInfoList"`
		Code    int         `json:"code"`
		Message interface{} `json:"message"`
	} `json:"privileges"`
	Code int `json:"code"`
}

func FetchPlaylist(playlistID string) (*[]player.Music, error) {
	logger := config.Logger
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.Config.NetEaseAPI+"/playlist/track/all", nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create request")
		return nil, err
	}
	req.Header.Set("Cookie", config.Config.NetEaseCookie)
	params := req.URL.Query()
	params.Add("id", playlistID)
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
		logger.Error().Msg("Failed to fetch playlist")
		return nil, err
	}
	var playlistResp PlaylistResp
	logger.Debug().Msg(string(body))
	if err := json.Unmarshal(body, &playlistResp); err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal response")
		return nil, err
	}
	var musicsList []player.Music
	for _, song := range playlistResp.Songs {
		musicResp, err := QueryMusic(song.ID)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get music info")
			continue
		}

		musicsList = append(musicsList, *musicResp)
	}
	return &musicsList, nil
}
