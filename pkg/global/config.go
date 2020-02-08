package global

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

type HttpConfig struct {
	Port      string `json:"port"`
	GzipLevel int    `json:"gzip_level"`
}

type DBConfig struct {
	Dsn     string `json:"dsn"`
	MaxOpen int    `json:"max_open"`
	MaxIdle int    `json:"max_idle"`
}

type AdminConfig struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type GlobalConfig struct {
	Http  *HttpConfig  `json:"http"`
	DB    *DBConfig    `json:"db"`
	Admin *AdminConfig `json:"admin"`
}

var (
	config *GlobalConfig
)

func Config() *GlobalConfig {
	return config
}

func ParseConfig() {
	if !file.IsExist(ConfigPath) {
		logrus.Fatalf("cannot find config file: %s", ConfigPath)
	}
	fullConfigPath, err := filepath.Abs(ConfigPath)
	if err != nil {
		logrus.Fatalf("cannot read config file: %s, error: %v", fullConfigPath, err)
	}
	configContent, err := file.ToTrimString(fullConfigPath)
	if err != nil {
		logrus.Fatalf("cannot read config file: %s, error: %v", fullConfigPath, err)
	}

	err = json.Unmarshal([]byte(configContent), &config)
	if err != nil {
		logrus.Fatalln("cannot parse config file: %s, error: %v", fullConfigPath, err)
	}
	logrus.Infoln("load config: " + fullConfigPath + " successfully")
}

func InitLog(level string) (err error) {
	switch level {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	default:
		logrus.Warn("log conf only allow [info, debug, warn]")
		return errors.New("log level set failed，only allow [info, debug, warn].")
	}
	logrus.Info("log level: " + level)
	return nil
}
