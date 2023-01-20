package config

import (
	"github.com/phuslu/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type envConfig struct {
	BotToken string `yaml:"BotToken"`
}

var Config *envConfig
var Logger *log.Logger

func LoadConfig(filename string) error {
	Logger = &log.Logger{
		Level:  log.InfoLevel,
		Writer: &log.ConsoleWriter{},
	}
	logger := Logger
	ymlFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read config file")
		return err
	}
	if err = yaml.Unmarshal(ymlFile, &Config); err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal config file")
		return err
	}
	return nil
}
