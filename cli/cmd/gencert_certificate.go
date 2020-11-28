package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

const (
	timeSkew = time.Duration(5) * time.Second // to protect against time-skew issues we potentially generate a certificate timeSkew duration
)

func genParseCertificate(bytes []byte) (*ssh.Certificate, error) {
	k, _, _, _, err := ssh.ParseAuthorizedKey(bytes)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Could not parse cert")
		return nil, fmt.Errorf("Could not parse cert: %w", err)
	}

	cert, ok := k.(*ssh.Certificate)
	if !ok {
		log.Warn().
			Err(err).
			Msg("Bytes do not correspond to an ssh certificate")
		return nil, fmt.Errorf("Bytes do not correspond to an ssh certificate: %w", err)
	}

	return cert, nil
}

func genCertFromFile() (*ssh.Certificate, error) {
	bytes, err := ioutil.ReadFile(genCertPath)
	if err != nil {
		log.Warn().
			Err(err).
			Str("filename", genCertPath).
			Msg("Could not read cert from file")
		return nil, fmt.Errorf("Could not read cert from file: %w", err)
	}

	return genParseCertificate(bytes)
}

func genIsCertStillFresh(cert *ssh.Certificate) bool {
	if cert == nil {
		return false
	}

	now := time.Now()
	validBefore := time.Unix(int64(cert.ValidBefore), 0).Add(-1 * timeSkew) // upper bound

	return now.Before(validBefore)
}

func genIsCertValidInCache() (bool, *ssh.Certificate) {
	cachedCert, err := genCertFromFile()

	if err == nil {
		isFresh := genIsCertStillFresh(cachedCert)

		if isFresh {
			log.Info().
				Str("certificate", genCertPath).
				Time("valid until", time.Unix(int64(cachedCert.ValidBefore), 0).UTC()).
				Msg("Current certificate still valid. Exiting...")
			return true, cachedCert
		}
	}

	return false, nil
}

func genStoreCertAtFile(filename string, cert []byte) error {
	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, cert, 0600)
}
