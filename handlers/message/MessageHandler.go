package message

import (
	"MusicBot/config"
	"MusicBot/middleware"
	"MusicBot/serve/Bili"
	"MusicBot/serve/LLM"
	"MusicBot/serve/NetEase"
	"MusicBot/serve/QQ"
	"MusicBot/serve/music"
	"MusicBot/serve/music/Functions"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/lonelyevil/kook"
)

func MessageHan(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	if ctx.Common.Type != kook.MessageTypeKMarkdown || ctx.Extra.Author.Bot {
		return
	}
	logger.Info().Msg("Message received: " + ctx.Common.Content)
	music.PlayStatus.Ctx = ctx
	if ctx.Common.Content == "/stop" {
		err := middleware.AdminMiddleware(ctx)
		if err != nil {
			return
		}

		music.PlayStatus.CanAppend = false
		music.SendMsg(ctx, "已停止添加新音乐")
		return
	} else if ctx.Common.Content == "/start" {
		err := middleware.AdminMiddleware(ctx)
		if err != nil {
			return
		}

		music.PlayStatus.CanAppend = true
		music.SendMsg(ctx, "可以添加新音乐")
		return
	}
	if strings.HasPrefix(ctx.Common.Content, "ping") {
		// Ping
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "ping")
		_, _ = ctx.Session.MessageCreate(&kook.MessageCreate{
			MessageCreateBase: kook.MessageCreateBase{
				TargetID: ctx.Common.TargetID,
				Content:  "pong",
				Quote:    ctx.Common.MsgID,
			},
		})
	} else if strings.HasPrefix(ctx.Common.Content, "/version") {
		// Version
		_, _ = ctx.Session.MessageCreate(&kook.MessageCreate{
			MessageCreateBase: kook.MessageCreateBase{
				TargetID: ctx.Common.TargetID,
				Content:  "golang 1.0",
				Quote:    ctx.Common.MsgID,
			},
		})
	} else if strings.HasPrefix(ctx.Common.Content, "/reload") {
		// Reload
		err := middleware.AdminMiddleware(ctx)
		if err != nil {
			return
		}

		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "ping")
		err = config.LoadConfig("config/config.yaml")
		if err != nil {
			logger.Error().Err(err).Msg("Reload config failed")
			return
		}
		_, _ = ctx.Session.MessageCreate(&kook.MessageCreate{
			MessageCreateBase: kook.MessageCreateBase{
				TargetID: ctx.Common.TargetID,
				Content:  "配置文件已重载",
				Quote:    ctx.Common.MsgID,
			},
		})
	} else if strings.HasPrefix(ctx.Common.Content, "/v ") {
		// Change volume
		err := middleware.AdminMiddleware(ctx)
		if err != nil {
			return
		}

		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/v ")
		ChangeVolumeMessageHandler(ctx)
	} else if strings.HasPrefix(ctx.Common.Content, "/c ") {
		// Change channel
		err := middleware.AdminMiddleware(ctx)
		if err != nil {
			return
		}

		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/c ")
		ChangeChannelMessageHandler(ctx)
	} else if strings.HasPrefix(ctx.Common.Content, "/n ") || strings.HasPrefix(ctx.Common.Content, "/s ") || strings.HasPrefix(ctx.Common.Content, "/b ") || strings.HasPrefix(ctx.Common.Content, "/q ") {
		// If command is from valid whitelist channel

		err := middleware.WhitelistChannelMiddleware(ctx)
		if err != nil {
			return
		}

		if strings.HasPrefix(ctx.Common.Content, "/n ") {
			// Netease Music
			ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/n ")
			NetEaseMusicMessageHandler(ctx)
		} else if strings.HasPrefix(ctx.Common.Content, "/s ") {
			// Netease Music Search
			ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/s ")
			NetEaseMusicSearchMessageHandler(ctx)
		} else if strings.HasPrefix(ctx.Common.Content, "/b ") {
			// Bilibili
			ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/b ")
			BiliMessageHandler(ctx)
		} else if strings.HasPrefix(ctx.Common.Content, "/q ") {
			// QQ Music
			ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/q ")
			QQMusicMessageHandler(ctx)
		}
	} else if ctx.Common.Content == "/skip" {
		// Skip
		SkipMusicMessageHandler(ctx)
	} else if ctx.Common.Content == "/list" {
		// List
		Functions.SendMusicList(ctx, &music.Musics)
	} else if strings.HasPrefix(ctx.Common.Content, "/llm ") {
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/llm ")

		resp, err := LLM.LLMQuery(ctx.Common.Content)
		if err != nil {
			music.SendMsg(ctx, fmt.Sprintf("Error: %s, Response: %s", err.Error(), resp))
		} else {
			music.SendMsg(ctx, resp)
		}
	}
}

func NetEaseMusicSearchMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	searchList, err := NetEase.SearchMusic(ctx.Common.Content)
	if err != nil {
		logger.Error().Err(err).Msg("Search music failed")
		music.SendMsg(ctx, "搜索音乐失败")
		return
	}
	NetEase.SendSelectList(ctx, searchList)
}

func NetEaseMusicMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger

	re := regexp.MustCompile(`song\?id=(\d+)`)
	match := re.FindStringSubmatch(ctx.Common.Content)
	if len(match) > 1 {
		ctx.Common.Content = match[1]
	}

	id, err := strconv.ParseInt(ctx.Common.Content, 10, 64)
	if err != nil {
		logger.Error().Err(err).Msg("Parse music id failed")
		music.SendMsg(ctx, "解析音乐ID失败：输入不合法")
		return
	}

	logger.Info().Msgf("Query music for id: %d", id)
	musicResult, err := NetEase.QueryMusic(int(id))
	if err != nil {
		logger.Error().Err(err).Msg("Query music failed")
		music.SendMsg(ctx, "查询音乐失败")
		return
	}
	music.SendMsg(ctx, fmt.Sprintf("%s 已加入播放列表", musicResult.Name))
	music.Musics.Add(musicResult)
	go music.Musics.Play(ctx)
}

func QQMusicMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	if strings.Contains(ctx.Common.Content, "songDetail/") {
		ctx.Common.Content = strings.Split(strings.Split(ctx.Common.Content, "songDetail/")[1], "]")[0]
	}
	id := ctx.Common.Content
	musicResult, err := QQ.QueryMusic(id)
	if err != nil {
		logger.Error().Err(err).Msg("Query music failed")
		music.SendMsg(ctx, "查询音乐失败")
		return
	}
	music.SendMsg(ctx, fmt.Sprintf("%s 已加入播放列表", musicResult.Name))
	music.Musics.Add(musicResult)
	go music.Musics.Play(ctx)
}

func BiliMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	if !(strings.HasPrefix(ctx.Common.Content, "av") || strings.HasPrefix(ctx.Common.Content, "AV") || strings.HasPrefix(ctx.Common.Content, "bv") || strings.HasPrefix(ctx.Common.Content, "BV")) {
		logger.Error().Msg("Parse video id failed")
		music.SendMsg(ctx, "解析AV/BV失败：输入不合法")
		return
	}
	var musicResult *music.Music
	var err error
	if strings.HasPrefix(ctx.Common.Content, "av") {
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "av")
		musicResult, err = Bili.QueryBiliAudio(ctx.Common.Content, false)
	} else if strings.HasPrefix(ctx.Common.Content, "AV") {
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "AV")
		musicResult, err = Bili.QueryBiliAudio(ctx.Common.Content, false)
	} else if strings.HasPrefix(ctx.Common.Content, "bv") {
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "bv")
		musicResult, err = Bili.QueryBiliAudio(ctx.Common.Content, true)
	} else if strings.HasPrefix(ctx.Common.Content, "BV") {
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "BV")
		musicResult, err = Bili.QueryBiliAudio(ctx.Common.Content, true)
	}
	if err != nil {
		logger.Error().Err(err).Msg("Query audio failed")
		music.SendMsg(ctx, "查询视频失败："+err.Error())
		return
	}
	music.SendMsg(ctx, fmt.Sprintf("%s 已加入播放列表", musicResult.Name))
	music.Musics.Add(musicResult)
	go music.Musics.Play(ctx)
}

func ChangeVolumeMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	if ctx.Common.Content == "now" {
		music.SendMsg(ctx, fmt.Sprintf("当前音量: %d", music.PlayStatus.Volume))
		return
	}
	if ctx.Common.Content == "anal" {
		if music.PlayStatus.Music == nil {
			music.SendMsg(ctx, "当前无音乐播放")
			return
		}
		analRespChan, logChan := music.StatVolume(music.PlayStatus.Music.File)
		analResp := <-analRespChan
		log := <-logChan
		msg := fmt.Sprintf("音量分析: 输入 %.2f(dB)/%.2f(LUFS) -> 音量均衡 %.2f(dB)/%.2f(LUFS) -> 增益 %.2f(dB)", analResp.OriginVolumedB, analResp.OriginVolumeLUFS, analResp.LoudnormVolumedB, analResp.LoudnormVolumeLUFS, analResp.GainVolumedB)
		logger.Info().Msg(log)
		logger.Info().Msg(msg)
		music.SendMsg(ctx, msg)
		return
	}
	volume, err := strconv.ParseInt(ctx.Common.Content, 10, 64)
	if err != nil {
		logger.Error().Err(err).Msg("Parse volume failed")
		music.SendMsg(ctx, "解析音量失败：输入不合法")
		return
	}
	Functions.ChangeVolume(int(volume))
	music.SendMsg(ctx, fmt.Sprintf("音量已调整为 %ddB", volume))
}

func ChangeChannelMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	if ctx.Common.Content == "list" {
		music.SendMsg(ctx, fmt.Sprintf("可用频道列表: %s", config.ListChannel()))
		return
	}
	channel, err := config.FindChannelID(ctx.Common.Content)
	if err != nil {
		logger.Error().Err(err).Msg("Find channel failed")
		music.SendMsg(ctx, "解析频道失败：频道不存在")
		return
	}
	Functions.ChangeChannel(channel)
	music.SendMsg(ctx, fmt.Sprintf("频道已切换为 %s", ctx.Common.Content))
}

func SkipMusicMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	if music.PlayStatus.Music == nil {
		music.SendMsg(ctx, "当前无音乐播放")
		return
	}
	logger.Info().Msg(fmt.Sprintf("Skipped music %s", music.PlayStatus.Music.Name))
	music.SendMsg(ctx, fmt.Sprintf("将跳过歌曲 %s", music.PlayStatus.Music.Name))
	Functions.SkipMusic()
}
