package main

import (
	"MusicBot/bff"
	"MusicBot/config"
	"MusicBot/handlers"
	"MusicBot/handlers/message"
	"MusicBot/serve/music"
	"github.com/gin-gonic/gin"
	"github.com/lonelyevil/kook"
	"github.com/lonelyevil/kook/log_adapter/plog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// Load config
	err := config.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err)
		return
	}
	// Setup logger
	logger := config.Logger
	// Setup KOOK
	s := kook.New(config.Config.BotToken, plog.NewLogger(logger))
	// Setup Gin
	ginServer := gin.Default()
	// Setup Env
	err = music.InitMusicEnv()
	if err != nil {
		logger.Error().Err(err).Msg("Init music env failed")
		return
	}

	// Register KOOK handlers
	handlers.RegistryHandlers(s, message.MessageHan)
	// Register Gin handlers
	bff.RegistryHandlers(ginServer)

	// Start KOOK
	err = s.Open()
	if err != nil {
		panic(err)
		return
	}
	logger.Info().Msg("Bot is running")
	// Start Gin
	err = ginServer.Run(":" + strconv.Itoa(config.Config.WebPort))
	if err != nil {
		panic(err)
		return
	}
	logger.Info().Msg("Web server is running")
	// Waiting for exit signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sc
	logger.Info().Msg("Bot is shutting down")
	err = s.Close()
}
