package cmd

import (
	"github.com/le0pard/certonid/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	radnstrCmd = &cobra.Command{
		Use:   "randstr",
		Short: "Generate random string (32 bytes)",
		Long:  `Generate random string (32 bytes), which can be used for CERTONID_SYMMETRIC_KEY environment variable`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err  error
				text string
			)

			text, err = utils.GenerateRandomString(32)

			if err != nil {
				er(err)
			}

			log.WithFields(log.Fields{
				"string": text,
			}).Info("Successfully generated random string")
		},
	}
)

func init() {
	rootCmd.AddCommand(radnstrCmd)
}
