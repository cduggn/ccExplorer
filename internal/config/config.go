package config

import "github.com/spf13/viper"

var LoadConfigFunc = func(path string) func() {
	return func() {
		LoadConfig(path)
	}
}

func LoadConfig(path string) {
	viper.AutomaticEnv()
}
