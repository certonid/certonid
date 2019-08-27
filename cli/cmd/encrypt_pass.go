package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(encryptPassCmd)
}

var encryptPassCmd = &cobra.Command{
	Use:   "encryptpass",
	Short: "Encrypt password",
	Long:  `Encrypt password with symmetric or kms encryption`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Certonid v0.1.0")
	},
}
