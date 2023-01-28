package Functions

import (
	"MusicBot/config"
	"MusicBot/serve/music"
	"fmt"
	"github.com/lonelyevil/kook"
	"strings"
	"time"
)

func SendMusicList(ctx *kook.KmarkdownMessageContext, musicsList *music.MusicsList) {
	logger := config.Logger
	cardMsg := kook.CardMessageCard{
		Theme: kook.CardThemeSuccess,
		Size:  kook.CardSizeLg,
	}
	for _, musicCtx := range musicsList.Musics {
		section := kook.CardMessageSection{
			Mode: kook.CardMessageSectionModeRight,
			Text: kook.CardMessageElementKMarkdown{
				Content: "**歌曲：** " + musicCtx.Name + "\n**歌手：** " + strings.Join(musicCtx.Artists, ", ") + "\n**时长：** " + time.Duration(musicCtx.LastTime*1000000).String(),
			},
		}
		cardMsg.AddModule(section.SetAccessory(&kook.CardMessageElementButton{
			Theme: kook.CardThemeDanger,
			Text:  "删除",
			Click: string(kook.CardMessageElementButtonClickReturnVal),
			Value: "DEL" + musicCtx.ID,
		}))
	}
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
