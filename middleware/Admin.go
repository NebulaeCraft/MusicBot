package middleware

import (
	"MusicBot/config"
	"MusicBot/serve/music"
	"MusicBot/serve/privilege"
	"errors"
	"github.com/lonelyevil/kook"
)

func AdminMiddleware(ctx *kook.KmarkdownMessageContext) error {
	logger := config.Logger

	isAdmin, err := privilege.IsAdmin(ctx.Common.AuthorID)
	if err != nil {
		logger.Error().Err(err).Msgf("IsAdmin function error: %s", ctx.Common.AuthorID)
		return errors.New("unable to find admin list")
	}
	if !isAdmin {
		music.SendMsg(ctx, "该命令仅管理员可用")
		logger.Error().Msg("insufficient privileges for userID: " + ctx.Common.AuthorID)
		return errors.New("insufficient privileges for userID: " + ctx.Common.AuthorID)
	}

	return nil
}
