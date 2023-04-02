package handlers

import (
	"fmt"
	handlers "github.com/cduggn/ccexplorer/internal/core/handlers/aws"
	"github.com/cduggn/ccexplorer/internal/core/logger"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ccexplorer",
		Short: "A CLI tool to explore cloud costs and usage",
		Long:  paintRootHeader(),
	}
	cfgFile        string
	LoadConfigFunc = func(path string) func() {
		return func() {
			LoadConfig(path)
		}
	}
)

func RootCommand() *cobra.Command {
	_, err := logger.New()
	if err != nil {
		panic(err.Error())
	}

	LoadConfigFunc(".")()
	handlers.Initialize()

	//cobra.OnInitialize(LoadConfigFunc("."))
	return rootCmd
}

func init() {
	rootCmd.AddCommand(handlers.CostAndForecast())
	rootCmd.AddCommand(handlers.Presets())
	_ = viper.BindPFlag("open_ai_api_key", rootCmd.PersistentFlags().Lookup(
		"OPEN_AI_API_KEY"))
	_ = viper.BindPFlag("aws_profile", rootCmd.PersistentFlags().Lookup(
		"AWS_PROFILE"))
}

func LoadConfig(path string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(path)
		viper.SetConfigType("env")
		viper.SetConfigName(".ccexplorer")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file specified:", err.Error())
	}
}

func paintRootHeader() string {
	myFigure := figure.NewFigure("ccExplorer", "thin", true)
	return myFigure.String()
}
