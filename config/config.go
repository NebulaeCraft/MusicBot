package config

import (
	"errors"
	"github.com/phuslu/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

type envConfig struct {
	BotToken      string `yaml:"BotToken"`
	WebPort       int    `yaml:"WebPort"`
	NetEaseAPI    string `yaml:"NetEaseApi"`
	NetEaseCookie string `yaml:"NetEaseCookie"`
	QQAPI         string `yaml:"QQApi"`
	VoicePort     int    `yaml:"VoicePort"`
	KOOKVoice     string `yaml:"KOOKVoice"`
	SearchLimit   int    `yaml:"SearchLimit"`
	VoiceChannel  []struct {
		Name string `yaml:"Name"`
		ID   int64  `yaml:"ID"`
	} `yaml:"VoiceChannel"`
	WhitelistChannel []int `yaml:"WhitelistChannel"`
	AdminUser        []struct {
		Name string `yaml:"Name"`
		ID   int    `yaml:"ID"`
	} `yaml:"AdminUser"`
	LLM struct {
		SiliconFlowToken string  `yaml:"SiliconFlowToken"`
		Model            string  `yaml:"Model"`
		Temperature      float64 `yaml:"Temperature"`
	} `yaml:"LLM"`
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

func FindChannelID(name string) (int64, error) {
	for _, v := range Config.VoiceChannel {
		if v.Name == name {
			return v.ID, nil
		}
	}
	return 0, errors.New("channel not found")
}

func ListChannel() string {
	var list []string
	for _, v := range Config.VoiceChannel {
		list = append(list, v.Name)
	}
	return strings.Join(list, ", ")
}
