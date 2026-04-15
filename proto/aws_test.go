package proto

import (
	"encoding/json"
	"testing"
)

func TestAwsSignEvent_JSON(t *testing.T) {
	event := AwsSignEvent{
		CertType:     "user",
		Key:          "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
		Username:     "testuser",
		Hostnames:    "",
		ValidUntil:   "1h",
		KMSAuthToken: "token123",
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded AwsSignEvent
	err = json.Unmarshal(bytes, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Username != "testuser" || decoded.KMSAuthToken != "token123" {
		t.Errorf("Decoded event does not match original: %+v", decoded)
	}
}

func TestAwsSignResponse_JSON(t *testing.T) {
	resp := AwsSignResponse{
		Cert: "ssh-rsa-cert-v01@openssh.com AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded AwsSignResponse
	err = json.Unmarshal(bytes, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.Cert != resp.Cert {
		t.Errorf("Decoded response does not match original: %+v", decoded)
	}
}
