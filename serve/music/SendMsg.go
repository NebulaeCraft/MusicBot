package music

import "github.com/lonelyevil/kook"

func SendMsg(ctx *kook.KmarkdownMessageContext, content string) {
	_, _ = ctx.Session.MessageCreate(&kook.MessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			TargetID: ctx.Common.TargetID,
			Content:  content,
		},
	})
}
