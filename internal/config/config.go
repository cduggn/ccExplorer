package config

import "github.com/spf13/viper"

// LoadConfigFunc is a function that loads the configuration
var LoadConfigFunc = func(path string) func() {
	return func() {
		LoadConfig(path)
	}
}

func LoadConfig(path string) {
	viper.AutomaticEnv()
}
