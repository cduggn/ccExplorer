package cmd

import (
	"github.com/cduggn/ccexplorer/internal/pkg/cmd/get"
	"github.com/cduggn/ccexplorer/internal/pkg/logger"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	//"go.uber.org/zap"
	//"go.uber.org/zap/zapcore"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ccexplorer",
		Short: "A CLI tool to explore cloud costs and usage",
		Long:  paintRootHeader(),
	}
)

func init() {
	rootCmd.AddCommand(get.AWSCostAndUsageCommand())
	rootCmd.AddCommand(AddPresetCommands())
}

func paintRootHeader() string {
	myFigure := figure.NewFigure("ccXplorer", "thin", true)
	return myFigure.String()
}

func Execute() {

	_, err := logger.New()
	if err != nil {
		panic(err.Error())
	}

	//defer logger.Sync()
	err = rootCmd.Execute()
	//rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	// Bind the flag to a Viper variable
	//viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))

	if err != nil {
		os.Exit(126)
	}
}

//func setLoggingLevel() zapcore.Level {
//	switch viper.GetString("log") {
//	case "debug":
//		return zap.DebugLevel
//	case "info":
//		return zap.InfoLevel
//	case "warning":
//		return zap.WarnLevel
//	case "error":
//		return zap.ErrorLevel
//	case "critical":
//		return zap.FatalLevel
//	}
//}
