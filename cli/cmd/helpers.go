package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func er(msg interface{}) {
	log.Error(msg)
	os.Exit(1)
}
