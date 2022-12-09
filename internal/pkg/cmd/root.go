package cmd

import (
	"github.com/cduggn/cloudcost/internal/pkg/storage"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cloudcost",
		Short: "A CLI tool to get AWS Costs, Usage and Forecasts",
		Long:  paintRootHeader(),
	}
	Database *storage.CostDataStorage
)

func init() {
	rootCmd.AddCommand(AWSCostAndUsageCommand())
}

func paintRootHeader() string {
	myFigure := figure.NewFigure("CloudCost", "thin", true)
	return myFigure.String()
}

func Execute(db *storage.CostDataStorage) {
	Database = db
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(126)
	}
}

//func New() {
//
//	database = storage.New("sqlite3", "./cloudcost.db")
//	rootCmd.AddCommand(billing.CostAndUsageCommand())
//
//	record := storage.CostDataInsert{
//		Dimension:   "test",
//		Dimension2:  "test2",
//		Tag:         "test3",
//		MetricName:  "test4",
//		Amount:      1.0,
//		Unit:        "test5",
//		Granularity: "test6",
//		StartDate:   "test7",
//		EndDate:     "test8",
//	}
//
//	storage.InsertCustomer(database, record)
//
//}
