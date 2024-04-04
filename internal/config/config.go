package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config - конфигурация приложения
type Config struct {
	BotId               string   `mapstructure:"bot_id"`                 // айди бота от ОтцаБотов
	EnablePhrases       bool     `mapstructure:"enable_phrases"`         // включить фиксированные фразы
	DefaultPhrases      []string `mapstructure:"default_phrases"`        // список фраз
	EnableSemen         bool     `mapstructure:"enable_semen"`           // включить генерацию фраз
	Ratio               int      `mapstructure:"ratio"`                  // количество сообщений между ответами бота
	Length              int      `mapstructure:"ratio"`                  // длина сообщений генератоа цепей
	DefaultDataFileName string   `mapstructure:"default_data_file_name"` // текстовый файл из которого берутся базовые данные
	MainCum             string   `mapstructure:"main_cum"`               // ник владельца
}

// LoadConfig - загружает конфигурацию приложения из указанного в path файла
func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	v := viper.New()
	v.SetConfigFile(path)

	err := v.ReadInConfig()

	if err != nil {
		return nil, errors.Wrap(err, "read config")
	}
	v.AutomaticEnv()

	err = v.Unmarshal(
		config,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.TextUnmarshallerHookFunc(),
				mapstructure.StringToTimeDurationHookFunc(),
			),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal config")
	}
	setDefaults(config)

	if err := Validate(config); err != nil {
		return nil, errors.Wrap(err, "validate config")
	}

	return config, nil
}

func setDefaults(c *Config) {
	c.Ratio = 50
	c.Length = 50
}

func (c *Config) Validate() error {
	if c.BotId == "" {
		return errors.New("Bot ID is missing")
	}

	return nil
}
