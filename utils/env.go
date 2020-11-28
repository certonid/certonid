package utils

import (
	"fmt"
	"os"
	"strings"
)

// EnvPrefix for environment variables
const EnvPrefix string = "CERTONID"

// EnvStrReplacer for environment variables
var EnvStrReplacer = strings.NewReplacer(".", "_")

// GetENV return env key as string
func GetENV(key string) (string, bool) {
	return os.LookupEnv(fmt.Sprintf("%s_%s", EnvPrefix, strings.ToUpper(key)))
}
