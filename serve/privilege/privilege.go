package privilege

import (
	"MusicBot/config"
	"strconv"
)

func IsInWhitelistChannel(channelID string) (bool, error) {
	id, err := strconv.ParseInt(channelID, 10, 64)
	if err != nil {
		return false, err
	}
	for _, v := range config.Config.WhitelistChannel {
		if v == int(id) {
			return true, nil
		}
	}
	return false, nil
}

func IsAdmin(authorID string) (bool, error) {
	id, err := strconv.ParseInt(authorID, 10, 64)
	if err != nil {
		return false, err
	}
	for _, v := range config.Config.AdminUser {
		if v.ID == int(id) {
			return true, nil
		}
	}
	return false, nil
}
