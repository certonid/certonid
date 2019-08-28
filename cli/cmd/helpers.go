package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func er(msg interface{}) {
	log.Error(msg)
	os.Exit(1)
}

func showHelp(cmd *cobra.Command, args []string) error {
	cmd.HelpFunc()(cmd, args)
	return nil
}
