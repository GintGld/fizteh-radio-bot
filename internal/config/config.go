package config

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	Log         Log    `yaml:"log"`
	RadioAddr   string `yaml:"radio-addr" env-required:"true"`
	WebhookAddr string `yaml:"webhook-addr" env-default:"8443"`
	TmpDir      string `yaml:"tmp-dir" env-default:"tmp"`
	UseFiller   bool   `yaml:"use-filler" env-default:"false"`
}

type Log struct {
	Srv Logger `yaml:"srv" env-default:""`
	Tg  Logger `yaml:"tg" env-default:""`
}

type Logger struct {
	Level  slog.Level `yaml:"level" env-default:"info"`
	Path   string     `yaml:"path" env-default:""`
	Pretty bool       `yaml:"pretty" env-default:"false"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
