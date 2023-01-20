package main

import (
	"MusicBot/config"
	"MusicBot/handlers"
	"MusicBot/handlers/message"
	"github.com/lonelyevil/kook"
	"github.com/lonelyevil/kook/log_adapter/plog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := config.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err)
		return
	}
	logger := config.Logger
	s := kook.New(config.Config.BotToken, plog.NewLogger(logger))
	err = s.Open()
	if err != nil {
		panic(err)
		return
	}

	handlers.RegistryHandlers(s, message.MessageHan)

	logger.Info().Msg("Bot is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sc
	logger.Info().Msg("Bot is shutting down")
	err = s.Close()
}
