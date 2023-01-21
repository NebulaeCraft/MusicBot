package message

import (
	"MusicBot/config"
	"MusicBot/serve/music"
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
	if strings.Contains(ctx.Common.Content, "ping") {
		id, _ := strconv.ParseInt(strings.TrimPrefix(ctx.Common.Content, "pingM"), 10, 64)
		music, err := music.QueryMusic(int(id))
		if err != nil {
			logger.Error().Err(err).Msg("Query music failed")
			return
		}
		_, _ = ctx.Session.MessageCreate(&kook.MessageCreate{
			MessageCreateBase: kook.MessageCreateBase{
				TargetID: ctx.Common.TargetID,
				Content:  fmt.Sprintf("%#v", music),
				Quote:    ctx.Common.MsgID,
				Type:     kook.MessageTypeKMarkdown,
			},
		})
	}
}
