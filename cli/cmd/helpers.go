package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func er(msg interface{}) {
	resError, ok := msg.(*error)

	if ok {
		log.Error().
			Err(*resError).
			Msg("Error")
	} else {
		log.Error().
			Msgf("Error %+v", msg)
	}

	os.Exit(1)
}

func showHelp(cmd *cobra.Command, args []string) error {
	cmd.HelpFunc()(cmd, args)
	return nil
}
