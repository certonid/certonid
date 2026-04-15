package utils

import "time"

const (
	// UserCertType mark as user certificate
	UserCertType string = "user"
	// HostCertType mark as host certificate
	HostCertType string = "host"
	// default timeSkew
	TimeSkew time.Duration = time.Duration(5) * time.Minute
)
