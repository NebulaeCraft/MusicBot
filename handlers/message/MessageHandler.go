package message

import (
	"MusicBot/config"
	"MusicBot/serve/music"
	"MusicBot/serve/music/bili"
	"fmt"
	"github.com/lonelyevil/kook"
	"strconv"
	"strings"
)

func MessageHan(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	if ctx.Common.Type != kook.MessageTypeKMarkdown || ctx.Extra.Author.Bot {
		return
	}
	logger.Info().Msg("Message received: " + ctx.Common.Content)
	if strings.HasPrefix(ctx.Common.Content, "/n ") {
		// Netease Music
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/n ")
		NeteaseMusicMessageHandler(ctx)
	} else if strings.HasPrefix(ctx.Common.Content, "ping") {
		// Ping
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "ping")
		_, _ = ctx.Session.MessageCreate(&kook.MessageCreate{
			MessageCreateBase: kook.MessageCreateBase{
				TargetID: ctx.Common.TargetID,
				Content:  "pong",
				Quote:    ctx.Common.MsgID,
			},
		})
	} else if strings.HasPrefix(ctx.Common.Content, "/v ") {
		// Change volume
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/v ")
		ChangeVolumeMessageHandler(ctx)
	} else if strings.HasPrefix(ctx.Common.Content, "/c ") {
		// Change channel
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/c ")
		ChangeChannelMessageHandler(ctx)
	} else if strings.HasPrefix(ctx.Common.Content, "/b ") {
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/b ")
		BiliMessageHandler(ctx)
	}
}

func NeteaseMusicMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	id, err := strconv.ParseInt(ctx.Common.Content, 10, 64)
	if err != nil {
		logger.Error().Err(err).Msg("Parse music id failed")
		SendMsg(ctx, "解析音乐ID失败：输入不合法")
		return
	}
	musicResult, err := music.QueryMusic(int(id))
	if err != nil {
		logger.Error().Err(err).Msg("Query music failed")
		SendMsg(ctx, "查询音乐失败")
		return
	}
	SendMusicCard(ctx, musicResult)
	music.PlayMusic(musicResult)
}

func BiliMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	if !(strings.HasPrefix(ctx.Common.Content, "av") || strings.HasPrefix(ctx.Common.Content, "AV") || strings.HasPrefix(ctx.Common.Content, "bv") || strings.HasPrefix(ctx.Common.Content, "BV")) {
		logger.Error().Msg("Parse video id failed")
		SendMsg(ctx, "解析AV/BV失败：输入不合法")
		return
	}
	musicResult, err := bili.QueryBiliAudio(ctx.Common.Content)
	if err != nil {
		logger.Error().Err(err).Msg("Query audio failed")
		SendMsg(ctx, "查询视频失败")
		return
	}
	SendMusicCard(ctx, musicResult)
	music.PlayMusic(musicResult)
}

func ChangeVolumeMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	volume, err := strconv.ParseInt(ctx.Common.Content, 10, 64)
	if err != nil {
		logger.Error().Err(err).Msg("Parse volume failed")
		SendMsg(ctx, "解析音量失败：输入不合法")
		return
	}
	music.ChangeVolume(int(volume))
	SendMsg(ctx, fmt.Sprintf("音量已调整为 %ddB", volume))
}

func ChangeChannelMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	channel, err := config.FindChannelID(ctx.Common.Content)
	if err != nil {
		logger.Error().Err(err).Msg("Find channel failed")
		SendMsg(ctx, "解析频道失败：频道不存在")
		return
	}
	music.ChangeChannel(channel)
	SendMsg(ctx, fmt.Sprintf("频道已切换为 %s", ctx.Common.Content))
}

func SendMsg(ctx *kook.KmarkdownMessageContext, content string) {
	_, _ = ctx.Session.MessageCreate(&kook.MessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			TargetID: ctx.Common.TargetID,
			Content:  content,
			Quote:    ctx.Common.MsgID,
		},
	})
}
