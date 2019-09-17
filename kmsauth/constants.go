package kmsauth

import (
	"time"
)

const (
	// TimeFormat  = "%Y%m%dT%H%M%SZ"
	TimeFormat = "20060102T150405Z"
	// timeSkew how much to compensate for time skew
	timeSkew = time.Duration(5) * time.Minute
	// TokenVersion1 is a token version
	TokenVersion1 = 1
	// TokenVersion2 is a token version
	TokenVersion2 = 2
)

// TokenVersion is a token version
type TokenVersion int
