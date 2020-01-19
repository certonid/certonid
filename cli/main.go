package main

import (
	"os"

	"github.com/certonid/certonid/cli/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	cmd.Execute()
}
