package music

import (
	"MusicBot/config"
	"fmt"
	"github.com/lonelyevil/kook"
	"strings"
	"time"
)

func SendMusicCard(ctx *kook.KmarkdownMessageContext, musicReq *Music) {
	logger := config.Logger
	cardMsg := kook.CardMessageCard{
		Theme: kook.CardThemeSuccess,
		Size:  kook.CardSizeLg,
	}
	section := kook.CardMessageSection{
		Mode: kook.CardMessageSectionModeRight,
		Text: kook.CardMessageElementKMarkdown{
			Content: "**歌曲：** " + musicReq.Name + "\n**歌手：** " + strings.Join(musicReq.Artists, ", ") + "\n**时长：** " + time.Duration(musicReq.LastTime*1000000).String(),
		},
	}
	cardMsg.AddModule(section.SetAccessory(&kook.CardMessageElementImage{
		Src:  musicReq.Album,
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
			Type:     kook.MessageTypeCard,
		},
	})
}
