package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	AppID         uint64         `yaml:"appID"`
	AccessToken   string         `yaml:"accessToken"`
	Sandbox       bool           `yaml:"sandbox"`
	ServerConfigs []ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Type string `yaml:"type"`
	Url  string `yaml:"url"`
}

var AppConf Config
var BaseApi string

func init() {
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(file, &AppConf); err != nil {
		panic(err)
	}
	if AppConf.Sandbox {
		BaseApi = "https://sandbox.api.sgroup.qq.com"
	} else {
		BaseApi = "https://api.sgroup.qq.com"
	}
}
