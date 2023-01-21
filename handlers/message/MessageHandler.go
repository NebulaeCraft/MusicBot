package message

import (
	"MusicBot/config"
	"MusicBot/serve/music"
	"github.com/lonelyevil/kook"
	"strconv"
	"strings"
)

func MessageHan(ctx *kook.KmarkdownMessageContext) {
	if ctx.Common.Type != kook.MessageTypeKMarkdown || ctx.Extra.Author.Bot {
		return
	}
	if strings.HasPrefix(ctx.Common.Content, "/n ") {
		ctx.Common.Content = strings.TrimPrefix(ctx.Common.Content, "/n ")
		MusicMessageHandler(ctx)
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
}
