package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	completionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `Generates bash completion scripts.

To load completion run

. <(certonid completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(certonid completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(os.Stdout)
		},
	}
)

func init() {
	rootCmd.AddCommand(completionCmd)
}
