package conf

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/sirupsen/logrus"
	"os"
)

var BaseConf Config

type LogConf struct {
	Dir   string
	Level string
}
type ServerConf struct {
	Port int32
}

type Config struct {
	DataDir string
	Log     LogConf
	Server  ServerConf
}

func InitConf() {
	confPath := os.Getenv("SERVER_CONFIG")
	if confPath == "" {
		confPath = "./config.toml"
	}
	if configFile, err := os.ReadFile(confPath); err != nil {
		panic("read conf error: " + err.Error())
	} else if err = toml.Unmarshal(configFile, &BaseConf); err != nil {
		panic("conf file unmarshal error: " + err.Error())
	}

	if level, err := logrus.ParseLevel(BaseConf.Log.Level); err == nil {
		logrus.SetLevel(level)
	}

	logrus.Printf("%+v", BaseConf)
}
