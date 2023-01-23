package Bili

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type VideoInfoResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Bvid      string `json:"bvid"`
		Aid       int    `json:"aid"`
		Videos    int    `json:"videos"`
		Tid       int    `json:"tid"`
		Tname     string `json:"tname"`
		Copyright int    `json:"copyright"`
		Pic       string `json:"pic"`
		Title     string `json:"title"`
		Pubdate   int    `json:"pubdate"`
		Ctime     int    `json:"ctime"`
		Desc      string `json:"desc"`
		DescV2    []struct {
			RawText string `json:"raw_text"`
			Type    int    `json:"type"`
			BizID   int    `json:"biz_id"`
		} `json:"desc_v2"`
		State     int `json:"state"`
		Duration  int `json:"duration"`
		MissionID int `json:"mission_id"`
		Rights    struct {
			Bp            int `json:"bp"`
			Elec          int `json:"elec"`
			Download      int `json:"download"`
			Movie         int `json:"movie"`
			Pay           int `json:"pay"`
			Hd5           int `json:"hd5"`
			NoReprint     int `json:"no_reprint"`
			Autoplay      int `json:"autoplay"`
			UgcPay        int `json:"ugc_pay"`
			IsCooperation int `json:"is_cooperation"`
			UgcPayPreview int `json:"ugc_pay_preview"`
			NoBackground  int `json:"no_background"`
			CleanMode     int `json:"clean_mode"`
			IsSteinGate   int `json:"is_stein_gate"`
			Is360         int `json:"is_360"`
			NoShare       int `json:"no_share"`
			ArcPay        int `json:"arc_pay"`
			FreeWatch     int `json:"free_watch"`
		} `json:"rights"`
		Owner struct {
			Mid  int    `json:"mid"`
			Name string `json:"name"`
			Face string `json:"face"`
		} `json:"owner"`
		Stat struct {
			Aid        int    `json:"aid"`
			View       int    `json:"view"`
			Danmaku    int    `json:"danmaku"`
			Reply      int    `json:"reply"`
			Favorite   int    `json:"favorite"`
			Coin       int    `json:"coin"`
			Share      int    `json:"share"`
			NowRank    int    `json:"now_rank"`
			HisRank    int    `json:"his_rank"`
			Like       int    `json:"like"`
			Dislike    int    `json:"dislike"`
			Evaluation string `json:"evaluation"`
			ArgueMsg   string `json:"argue_msg"`
		} `json:"stat"`
		Dynamic   string `json:"dynamic"`
		Cid       int    `json:"cid"`
		Dimension struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			Rotate int `json:"rotate"`
		} `json:"dimension"`
		Premiere           interface{} `json:"premiere"`
		TeenageMode        int         `json:"teenage_mode"`
		IsChargeableSeason bool        `json:"is_chargeable_season"`
		IsStory            bool        `json:"is_story"`
		NoCache            bool        `json:"no_cache"`
		Pages              []struct {
			Cid       int    `json:"cid"`
			Page      int    `json:"page"`
			From      string `json:"from"`
			Part      string `json:"part"`
			Duration  int    `json:"duration"`
			Vid       string `json:"vid"`
			Weblink   string `json:"weblink"`
			Dimension struct {
				Width  int `json:"width"`
				Height int `json:"height"`
				Rotate int `json:"rotate"`
			} `json:"dimension"`
			FirstFrame string `json:"first_frame"`
		} `json:"pages"`
		Subtitle struct {
			AllowSubmit bool `json:"allow_submit"`
			List        []struct {
				ID          int64  `json:"id"`
				Lan         string `json:"lan"`
				LanDoc      string `json:"lan_doc"`
				IsLock      bool   `json:"is_lock"`
				SubtitleURL string `json:"subtitle_url"`
				Type        int    `json:"type"`
				IDStr       string `json:"id_str"`
				AiType      int    `json:"ai_type"`
				AiStatus    int    `json:"ai_status"`
				Author      struct {
					Mid            int    `json:"mid"`
					Name           string `json:"name"`
					Sex            string `json:"sex"`
					Face           string `json:"face"`
					Sign           string `json:"sign"`
					Rank           int    `json:"rank"`
					Birthday       int    `json:"birthday"`
					IsFakeAccount  int    `json:"is_fake_account"`
					IsDeleted      int    `json:"is_deleted"`
					InRegAudit     int    `json:"in_reg_audit"`
					IsSeniorMember int    `json:"is_senior_member"`
				} `json:"author"`
			} `json:"list"`
		} `json:"subtitle"`
		IsSeasonDisplay bool `json:"is_season_display"`
		UserGarb        struct {
			URLImageAniCut string `json:"url_image_ani_cut"`
		} `json:"user_garb"`
		HonorReply struct {
		} `json:"honor_reply"`
		LikeIcon   string `json:"like_icon"`
		NeedJumpBv bool   `json:"need_jump_bv"`
	} `json:"data"`
}

type VideoInfo struct {
	Cover    string
	Up       string
	Title    string
	Duration int
}

func QueryVideoInfo(serial string, isBV bool) (*VideoInfo, error) {
	url := "https://api.bilibili.com/x/web-interface/view"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	params := req.URL.Query()
	if isBV {
		params.Add("bvid", serial)
	} else {
		params.Add("aid", serial)
	}
	req.URL.RawQuery = params.Encode()

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var videoInfo VideoInfoResp
	err = json.Unmarshal(body, &videoInfo)
	if err != nil {
		return nil, err
	}
	if videoInfo.Code != 0 {
		return nil, errors.New("code not 0")
	}
	return &VideoInfo{
		Cover:    videoInfo.Data.Pic + "@130w_130h.jpg",
		Up:       videoInfo.Data.Owner.Name,
		Title:    strings.Replace(videoInfo.Data.Title, "/", " ", -1),
		Duration: videoInfo.Data.Duration,
	}, nil
}
