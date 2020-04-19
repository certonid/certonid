package proto

// GcpSignRequest used for function arguments
type GcpSignRequest struct {
	CertType     string `json:"cert_type"`
	Key          string `json:"key"`
	Username     string `json:"username"`
	Hostnames    string `json:"hostnames"`
	ValidUntil   string `json:"valid_until"`
	KMSAuthToken string `json:"kmsauth_token"`
}

// GcpSignResponse used for function response
type GcpSignResponse struct {
	Cert string `json:"cert"`
}
