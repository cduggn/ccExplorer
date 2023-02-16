package commands

import (
	"fmt"
	aws_presets "github.com/cduggn/ccexplorer/internal/commands/aws-presets"
	"github.com/cduggn/ccexplorer/internal/commands/get"
	"github.com/cduggn/ccexplorer/pkg/logger"
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
	cfgFile string
)

func RootCommand() *cobra.Command {
	_, err := logger.New()
	if err != nil {
		panic(err.Error())
	}
	cobra.OnInitialize(initConfig)
	return rootCmd
}

func init() {
	rootCmd.AddCommand(get.AWSCostAndUsageCommand())
	rootCmd.AddCommand(aws_presets.AddAWSPresetCommands())
	err := viper.BindPFlag("open_ai_api_key", rootCmd.PersistentFlags().Lookup(
		"OPEN_AI_API_KEY"))
	if err != nil {
		fmt.Println("OPEN_AI_API_KEY not set")
	}
}

func paintRootHeader() string {
	myFigure := figure.NewFigure("ccExplorer", "thin", true)
	return myFigure.String()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file found:", err.Error())
	}
}
