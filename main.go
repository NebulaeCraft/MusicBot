package main

import (
	"MusicBot/config"
	"MusicBot/handlers"
	"MusicBot/handlers/button"
	"MusicBot/handlers/message"
	"MusicBot/serve/NetEase"
	"MusicBot/serve/player"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/lonelyevil/kook"
	"github.com/lonelyevil/kook/log_adapter/plog"
)

func main() {
	// Load config
	err := config.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}
	// Setup logger
	logger := config.Logger

	// Setup KOOK
	s := kook.New(config.Config.BotToken, plog.NewLogger(logger))

	// Setup Gin
	ginServer := gin.Default()

	// Setup Env
	player.MusicPlayer = player.NewPlayer()
	player.MusicPlayer.DefaultPlaylist, _ = NetEase.FetchPlaylist(config.Config.NetEaseDefaultPlaylist)
	logger.Info().Msgf("Default playlist length: %d", len(*player.MusicPlayer.DefaultPlaylist))

	// Register KOOK handlers
	handlers.RegistryHandlers(s, message.MessageHan, button.ButtonHan)

	// Register Gin handlers
	//bff.RegistryHandlers(ginServer)

	// Start KOOK
	err = s.Open()
	if err != nil {
		panic(err)
	}
	logger.Info().Msg("Bot is running")

	// Start Gin
	err = ginServer.Run(":" + strconv.Itoa(config.Config.WebPort))
	if err != nil {
		panic(err)
	}
	logger.Info().Msg("Web server is running")

	// Waiting for exit signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sc
	logger.Info().Msg("Bot is shutting down")
	err = s.Close()
}
