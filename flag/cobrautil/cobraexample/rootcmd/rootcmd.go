package rootcmd

import (
	"log/slog"
	"os"

	"github.com/grokify/sogo/flag/cobrautil/cobraexample"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cobraexample",
	Short: "Cobra example command",
	Long:  `Cobra example command test app`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	if cmd, err := cobraexample.ExampleCommand("testtypes"); err != nil {
		slog.Error("Error marking flag required", "errorMessage", err.Error())
		os.Exit(1)
	} else {
		rootCmd.AddCommand(cmd)
	}
}
