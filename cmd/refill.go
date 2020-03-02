package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/service/organizations"

	"github.com/giorgioazzinnaro/aws-poolboy/pkg/refill"
)

var refillCmd = &cobra.Command{
	Use:   "refill",
	Short: "Refill the targeted pool of accounts",
	Run:   refillRun,
}

var (
	argAccountsCount int
)

const (
	flagAccountsCount = "count"
)

func init() {
	refillCmd.PersistentFlags().IntVarP(
		&argAccountsCount,
		flagAccountsCount,
		"c",
		1,
		"The number of accounts to be created",
	)
}

func refillRun(cmd *cobra.Command, args []string) {

	orgAPI := organizations.New(sess)

	r := refill.New(
		orgAPI,
		&refill.RefillerOpts{
			Parallelism:         5,
			AccountPrefix:       "",
			AccountRoleName:     "",
			AccountRootDomain:   "",
			AccountRootUsername: "",
			CleanupOU:           "",
			TargetOU:            "",
		},
	)

	r.Create(argAccountsCount)
}
