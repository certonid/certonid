package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var gencertCmd = &cobra.Command{
	Use:   "gencert",
	Short: "Generate user or host certificate",
	Long:  `Generate user or host sertificate by involke serverless function`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}

func init() {
	rootCmd.AddCommand(gencertCmd)
}
