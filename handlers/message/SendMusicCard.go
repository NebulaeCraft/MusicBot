package message

import (
	"MusicBot/config"
	"MusicBot/serve/music"
	"fmt"
	"github.com/lonelyevil/kook"
	"strings"
	"time"
)

func SendMusicCard(ctx *kook.KmarkdownMessageContext, music *music.Music) {
	logger := config.Logger
	cardMsg := kook.CardMessageCard{
		Theme: kook.CardThemeSuccess,
		Size:  kook.CardSizeLg,
	}
	section := kook.CardMessageSection{
		Mode: kook.CardMessageSectionModeRight,
		Text: kook.CardMessageElementKMarkdown{
			Content: "**歌曲：** " + music.Name + "\n**歌手：** " + strings.Join(music.Artists, ", ") + "\n**时长：** " + time.Duration(music.LastTime*1000000).String(),
		},
	}
	cardMsg.AddModule(section.SetAccessory(&kook.CardMessageElementImage{
		Src:  music.Album,
		Size: "lg",
	}))
	cardMsgCtx, err := cardMsg.MarshalJSON()
	if err != nil {
		logger.Error().Err(err).Msg("Marshal card message failed")
		return
	}
	cardMsgCtxStr := fmt.Sprintf("[%s]", cardMsgCtx)
	_, _ = ctx.Session.MessageCreate(&kook.MessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			TargetID: ctx.Common.TargetID,
			Content:  cardMsgCtxStr,
			Quote:    ctx.Common.MsgID,
			Type:     kook.MessageTypeCard,
		},
	})
}
