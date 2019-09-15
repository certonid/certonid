package proto

// AwsSignEvent used for function arguments
type AwsSignEvent struct {
	CertType     string `json:"cert_type"`
	Key          string `json:"key"`
	Username     string `json:"username"`
	Hostnames    string `json:"hostnames"`
	ValidUntil   string `json:"valid_until"`
	KMSAuthToken string `json:"kmsauth_token"`
}

// AwsSignResponse used for function response
type AwsSignResponse struct {
	Cert string `json:"cert"`
}
