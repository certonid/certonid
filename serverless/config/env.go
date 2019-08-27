package config

import (
	"fmt"
	"os"
	"strings"
)

// EnvPrefix for environment variables
const EnvPrefix string = "CERTONID"

// GetENV return env key as string
func GetENV(key string) (string, bool) {
	return os.LookupEnv(fmt.Sprintf("%s_%s", EnvPrefix, strings.ToUpper(key)))
}
