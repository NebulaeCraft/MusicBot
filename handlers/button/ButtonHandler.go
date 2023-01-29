package button

import (
	"MusicBot/config"
	"MusicBot/serve/NetEase"
	"MusicBot/serve/music"
	"fmt"
	"github.com/lonelyevil/kook"
	"strconv"
	"strings"
)

func ButtonHan(ctx *kook.MessageButtonClickContext) {
	config.Logger.Info().Msg("Button received: " + ctx.Extra.Value)
	if strings.HasPrefix(ctx.Extra.Value, "NS") {
		ctx.Extra.Value = strings.TrimPrefix(ctx.Extra.Value, "NS")
		NetEaseSearchButtonHan(ctx)
	} else if strings.HasPrefix(ctx.Extra.Value, "DEL") {
		ctx.Extra.Value = strings.TrimPrefix(ctx.Extra.Value, "DEL")
		DeleteMusicButtonHan(ctx)
	} else if ctx.Extra.Value == "CONFIRM" {
		music.SendMsg(music.PlayStatus.Ctx, "你知道个🔨")
	}
}

func NetEaseSearchButtonHan(ctx *kook.MessageButtonClickContext) {
	logger := config.Logger
	if music.PlayStatus.CanAppend == false {
		music.SendMsg(music.PlayStatus.Ctx, "当前播放列表已锁定，无法添加新的音乐")
		return
	}
	id, _ := strconv.ParseInt(ctx.Extra.Value, 10, 64)
	musicResult, err := NetEase.QueryMusic(int(id))
	if err != nil {
		logger.Error().Err(err).Msg("Query music failed")
		music.SendMsg(music.PlayStatus.Ctx, "查询音乐失败")
		return
	}
	music.SendMsg(music.PlayStatus.Ctx, fmt.Sprintf("%s 已加入播放列表", musicResult.Name))
	music.Musics.Add(musicResult)
	go music.Musics.PlayBtn(ctx)
}

func DeleteMusicButtonHan(ctx *kook.MessageButtonClickContext) {
	logger := config.Logger
	index := music.Musics.GetIndexByID(ctx.Extra.Value)
	if index == -1 {
		logger.Error().Msg("Delete music failed")
		music.SendMsg(music.PlayStatus.Ctx, "删除音乐失败")
		return
	}
	music.SendMsg(music.PlayStatus.Ctx, fmt.Sprintf("已删除音乐 %s", music.Musics.Musics[index].Name))
	music.Musics.Musics = append(music.Musics.Musics[:index], music.Musics.Musics[index+1:]...)
}
