package message

import (
	"github.com/lonelyevil/kook"
	"strings"
)

func MessageHan(ctx *kook.KmarkdownMessageContext) {
	if ctx.Common.Type != kook.MessageTypeKMarkdown || ctx.Extra.Author.Bot {
		return
	}
	if strings.Contains(ctx.Common.Content, "ping") {
		_, _ = ctx.Session.MessageCreate(&kook.MessageCreate{
			MessageCreateBase: kook.MessageCreateBase{
				TargetID: ctx.Common.TargetID,
				Content:  "pong",
				Quote:    ctx.Common.MsgID,
				Type:     kook.MessageTypeKMarkdown,
			},
		})
	}
}
