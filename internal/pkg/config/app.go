package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"regexp"
)

type Config struct {
	ConsulHost  string         `yaml:"consulHost"`
	ConsulPort  int            `yaml:"consulPort"`
	ServiceName string         `yaml:"serviceName"`
	ServiceId   string         `yaml:"serviceId"`
	GrpcPort    int            `yaml:"grpcPort"`
	Bot         []BotConfig    `yaml:"bot"`
	Server      []ServerConfig `yaml:"server"`
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

	resolveEnv(&AppConf)

	for i := range AppConf.Bot {
		if AppConf.Bot[i].Sandbox {
			AppConf.Bot[i].Endpoint = "https://sandbox.api.sgroup.qq.com"
		} else {
			AppConf.Bot[i].Endpoint = "https://api.sgroup.qq.com"
		}
	}
}

func resolveEnv(cfg *Config) {
	cfgV := reflect.ValueOf(cfg).Elem()
	resolveStruct(cfgV)
}

func resolveStruct(v reflect.Value) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		switch field.Kind() {
		case reflect.String:
			resolveString(field, fieldType.Name)
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				resolveStruct(field.Index(j))
			}
		case reflect.Struct:
			resolveStruct(field)
		default:
		}
	}
}

func resolveString(v reflect.Value, fieldName string) {
	val := v.String()
	re := regexp.MustCompile(`\$\{(.+?)(?::(.+?))?}`)
	matches := re.FindStringSubmatch(val)

	if len(matches) > 1 {
		envKey := matches[1]
		defaultValue := ""
		if len(matches) == 3 {
			defaultValue = matches[2]
		}

		envVal, found := os.LookupEnv(envKey)
		if !found {
			if defaultValue == "" {
				panic(fmt.Sprintf("Environment variable for %s (%s) is not set and no default value is provided", fieldName, envKey))
			}
			envVal = defaultValue
		}
		v.SetString(envVal)
	}
}
