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
	cobra.OnInitialize(LoadConfigFunc("."))
	return rootCmd
}

func init() {
	rootCmd.AddCommand(get.AWSCostAndUsageCommand())
	rootCmd.AddCommand(aws_presets.AddAWSPresetCommands())
	_ = viper.BindPFlag("open_ai_api_key", rootCmd.PersistentFlags().Lookup(
		"OPEN_AI_API_KEY"))

}

func paintRootHeader() string {
	myFigure := figure.NewFigure("ccExplorer", "thin", true)
	return myFigure.String()
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
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file at:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file specified:", err.Error())
	}
}
