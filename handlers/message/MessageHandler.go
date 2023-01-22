package message

import (
	"MusicBot/config"
	"MusicBot/serve/music"
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
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/n ")
		MusicMessageHandler(ctx)
	} else if strings.HasPrefix(ctx.Common.Content, "ping") {
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "ping")
		_, _ = ctx.Session.MessageCreate(&kook.MessageCreate{
			MessageCreateBase: kook.MessageCreateBase{
				TargetID: ctx.Common.TargetID,
				Content:  "pong",
				Quote:    ctx.Common.MsgID,
			},
		})
	}
}

func MusicMessageHandler(ctx *kook.KmarkdownMessageContext) {
	logger := config.Logger
	id, _ := strconv.ParseInt(ctx.Common.Content, 10, 64)
	musicResult, err := music.QueryMusic(int(id))
	if err != nil {
		logger.Error().Err(err).Msg("Query music failed")
		return
	}
	SendMusicCard(ctx, musicResult)
	music.PlayMusic(musicResult)
}
