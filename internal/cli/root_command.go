package cli_new

import (
	"github.com/cduggn/ccexplorer/internal/config"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:     "ccexplorer",
		Version: "0.7.2",
		Short:   "A CLI tool to explore cloud costs and usage",
		Long:    paintRootHeader(),
	}
)

func RootCommand() *cobra.Command {
	config_new.LoadConfigFunc(".")()
	Initialize()
	return rootCmd
}

func init() {
	rootCmd.AddCommand(CostAndForecast())
	_ = viper.BindPFlag("openai_api_key", rootCmd.PersistentFlags().Lookup(
		"OPENAI_API_KEY"))
	_ = viper.BindPFlag("aws_profile", rootCmd.PersistentFlags().Lookup(
		"AWS_PROFILE"))
	_ = viper.BindPFlag("PINECONE_INDEX", rootCmd.PersistentFlags().Lookup(
		"PINECONE_INDEX"))
	_ = viper.BindPFlag("PINECONE_API_KEY", rootCmd.PersistentFlags().Lookup(
		"PINECONE_API_KEY"))
}

func paintRootHeader() string {
	myFigure := figure.NewFigure("ccExplorer", "thin", true)
	return myFigure.String()
}
