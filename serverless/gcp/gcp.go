package gcp

import (
	"encoding/json"
	"net/http"

	"github.com/certonid/certonid/proto"
	"github.com/certonid/certonid/serverless/config"
	"github.com/certonid/certonid/serverless/sshca"
	log "github.com/sirupsen/logrus"
)

// HandleSign function to generate certificate
func HandleSign(w http.ResponseWriter, r *http.Request) {
	var (
		requestData proto.GcpSignRequest
		err         error
		cert        string
	)

	config.InitConfig()

	jsonErr := json.NewDecoder(r.Body).Decode(&requestData)
	if jsonErr != nil {
		log.WithFields(log.Fields{
			"error": jsonErr,
		}).Error("Error parsing application/json")

		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}

	cert, err = sshca.GenerateCetrificate(&sshca.CertificateRequest{
		CertType:     requestData.CertType,
		Key:          requestData.Key,
		Username:     requestData.Username,
		Hostnames:    requestData.Hostnames,
		ValidUntil:   requestData.ValidUntil,
		KMSAuthToken: requestData.KMSAuthToken,
	})

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error to generate certificate")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	responseData, err := json.Marshal(proto.GcpSignResponse{Cert: cert})
	if err != nil {
		log.WithFields(log.Fields{
			"error": jsonErr,
		}).Error("Error marshal response in json")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}
