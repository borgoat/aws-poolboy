package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws/session"
)

var rootCmd = &cobra.Command{
	Use:   "poolboy",
	Short: "Poolboy takes care of AWS accounts",
}

var (
	sess *session.Session
)

func init() {

	cobra.OnInitialize(
		initConfig,
		initAws,
	)

	rootCmd.AddCommand(
		refillCmd,
	)
}

func initConfig() {

}

func initAws() {
	sess = session.Must(session.NewSession())
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
