package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Bot    []BotConfig    `yaml:"bot"`
	Server []ServerConfig `yaml:"server"`
}

type BotConfig struct {
	AppID       uint64 `yaml:"appID"`
	AccessToken string `yaml:"accessToken"`
	Sandbox     bool   `yaml:"sandbox"`
	Endpoint    string `yaml:"endpoint"`
	Tag         string `yaml:"tag"`
}

type ServerConfig struct {
	Type string `yaml:"type"`
	Url  string `yaml:"url"`
}

var AppConf Config

func init() {
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(file, &AppConf); err != nil {
		panic(err)
	}
	for i := range AppConf.Bot {
		if AppConf.Bot[i].Sandbox {
			AppConf.Bot[i].Endpoint = "https://sandbox.api.sgroup.qq.com"
		} else {
			AppConf.Bot[i].Endpoint = "https://api.sgroup.qq.com"
		}
	}
}
