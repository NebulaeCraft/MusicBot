package middleware

import (
	"MusicBot/config"
	"MusicBot/serve/music"
	"MusicBot/serve/privilege"
	"errors"
	"github.com/lonelyevil/kook"
)

func WhitelistChannelMiddleware(ctx *kook.KmarkdownMessageContext) error {
	logger := config.Logger

	isInWhitelistChannel, err := privilege.IsInWhitelistChannel(ctx.Common.TargetID)
	if err != nil {
		logger.Error().Err(err).Msgf("IsInWhitelistChannel function error: %s", ctx.Common.TargetID)
		return errors.New("unable to find whitelisted channel")
	}
	if !isInWhitelistChannel {
		music.SendMsg(ctx, "该频道不可点歌")
		logger.Error().Msg("channel not in whitelist for channel: " + ctx.Common.TargetID)
		return errors.New("channel not in whitelist")
	}

	return nil
}
