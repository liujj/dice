package g

import (
	"encoding/json"
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/toolkits/file"
	"sync"
)

type HttpConfig struct {
	AdminPort  string `json:"admin_port"`
	PublicPort string `json:"public_port"`
	GzipLevel  int    `json:"gzip_level"`
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
	Debug bool         `json:"debug"`
	Http  *HttpConfig  `json:"http"`
	DB    *DBConfig    `json:"db"`
	Log   string       `json:"log"`
	View  string       `json:"view"`
	Admin *AdminConfig `json:"admin"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		logrus.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		logrus.Fatalln("config file: ", cfg, " is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		logrus.Fatalln("read config file: ", cfg, " fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		logrus.Fatalln("parse config file: ", cfg, " fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c
	if c.Debug {
		InitLog("debug")
	} else {
		InitLog("warn")
	}
	logrus.Infoln("read config file: " + cfg + " successfully")
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
