package button

import (
	"MusicBot/config"
	"MusicBot/serve/NetEase"
	"MusicBot/serve/player"
	"fmt"
	"strconv"
	"strings"

	"github.com/lonelyevil/kook"
)

func ButtonHan(ctx *kook.MessageButtonClickContext) {
	config.Logger.Info().Msg("Button received: " + ctx.Extra.Value)
	if strings.HasPrefix(ctx.Extra.Value, "NS") {
		ctx.Extra.Value = strings.TrimPrefix(ctx.Extra.Value, "NS")
		NetEaseSearchButtonHan(ctx)
	} else if strings.HasPrefix(ctx.Extra.Value, "DEL") {
		ctx.Extra.Value = strings.TrimPrefix(ctx.Extra.Value, "DEL")
		DeleteMusicButtonHan(ctx)
	} else if strings.HasPrefix(ctx.Extra.Value, "UP") {
		ctx.Extra.Value = strings.TrimPrefix(ctx.Extra.Value, "UP")
		MoveUpMusicButtonHan(ctx)
	} else if strings.HasPrefix(ctx.Extra.Value, "DOWN") {
		ctx.Extra.Value = strings.TrimPrefix(ctx.Extra.Value, "DOWN")
		MoveDownMusicButtonHan(ctx)
	} else if strings.HasPrefix(ctx.Extra.Value, "TOP") {
		ctx.Extra.Value = strings.TrimPrefix(ctx.Extra.Value, "TOP")
		MoveTopMusicButtonHan(ctx)
	} else if ctx.Extra.Value == "CONFIRM" {
		player.MusicPlayer.SendMsg("你知道个🔨")
	}
}

func NetEaseSearchButtonHan(ctx *kook.MessageButtonClickContext) {
	logger := config.Logger
	id, _ := strconv.ParseInt(ctx.Extra.Value, 10, 64)
	musicResult, err := NetEase.QueryMusic(int(id))
	if err != nil {
		logger.Error().Err(err).Msg("Query music failed")
		player.MusicPlayer.SendMsg("查询音乐失败")
		return
	}
	player.MusicPlayer.AddMusic(musicResult)
}

func DeleteMusicButtonHan(ctx *kook.MessageButtonClickContext) {
	logger := config.Logger
	name := player.MusicPlayer.RemoveMusic(ctx.Extra.Value)
	if name == "" {
		logger.Error().Msg("Delete music failed")
		player.MusicPlayer.SendMsg("删除音乐失败")
		return
	}
	player.MusicPlayer.SendMsg(fmt.Sprintf("已删除音乐 %s", name))
}

func MoveUpMusicButtonHan(ctx *kook.MessageButtonClickContext) {
	logger := config.Logger
	name := player.MusicPlayer.MoveUpMusic(ctx.Extra.Value)
	if name == "" {
		logger.Error().Msg("Move up music failed")
		player.MusicPlayer.SendMsg("上移音乐失败")
		return
	}
	player.MusicPlayer.SendMsg(fmt.Sprintf("已上移音乐 %s", name))
}

func MoveDownMusicButtonHan(ctx *kook.MessageButtonClickContext) {
	logger := config.Logger
	name := player.MusicPlayer.MoveDownMusic(ctx.Extra.Value)
	if name == "" {
		logger.Error().Msg("Move down music failed")
		player.MusicPlayer.SendMsg("下移音乐失败")
		return
	}
	player.MusicPlayer.SendMsg(fmt.Sprintf("已下移音乐 %s", name))
}

func MoveTopMusicButtonHan(ctx *kook.MessageButtonClickContext) {
	logger := config.Logger
	name := player.MusicPlayer.MoveTopMusic(ctx.Extra.Value)
	if name == "" {
		logger.Error().Msg("Move top music failed")
		player.MusicPlayer.SendMsg("置顶音乐失败")
		return
	}
	player.MusicPlayer.SendMsg(fmt.Sprintf("已置顶音乐 %s", name))
}
