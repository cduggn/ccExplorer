package cmd

import (
	"github.com/cduggn/cloudcost/internal/pkg/cmd/billing"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "awsbill",
	Short: "A CLI tool to get AWS billing information",
	Long: `
--------------------------------------------
## # #   #  #    #    #  ##   #   # #
#     #  #  #    #    #  # #  #  # 
# # # #  #  #    #    #  #  # #  #    ##
#     #  #  #    #    #  #   ##  #    #
## # #   #  ###  ###  #  #    #   # # #
--------------------------------------------`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(126)
	}
}

func init() {
	rootCmd.AddCommand(billing.BillingCmd())
}
