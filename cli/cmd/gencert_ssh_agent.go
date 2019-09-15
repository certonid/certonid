package cmd

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/ScaleFT/sshkeys"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

func genGetPrivateKeyPassphrase(data []byte) ([]byte, error) {
	var (
		isEncryptedByPass bool
		passPhrase        []byte
		passPhraseErr     error
	)

	block, _ := pem.Decode(data)

	if block != nil {
		isEncryptedByPass = x509.IsEncryptedPEMBlock(block)

		if !isEncryptedByPass {
			_, err := x509.DecryptPEMBlock(block, []byte{})
			if err != nil {
				isEncryptedByPass = true
			}
		}

		if isEncryptedByPass {
			fmt.Print("SSH Key Passphrase [none]: ")
			passPhrase, passPhraseErr = terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Print("\n")
			if passPhraseErr != nil {
				log.WithFields(log.Fields{
					"error": passPhraseErr,
				}).Error("Not provided passphrase for provate key")
				return nil, passPhraseErr
			}
		}
	}

	return passPhrase, nil
}

func genCastInterfaceToPrimaryKeyInterface(key interface{}) interface{} {
	rsaKey, ok := key.(rsa.PrivateKey)
	if ok {
		return &rsaKey
	}
	dsaKey, ok := key.(dsa.PrivateKey)
	if ok {
		return &dsaKey
	}
	ecdsaKey, ok := key.(ecdsa.PrivateKey)
	if ok {
		return &ecdsaKey
	}
	ed25519Key, ok := key.(ed25519.PrivateKey)
	if ok {
		return &ed25519Key
	}
	return nil
}

func genAddCertToAgent(cert *ssh.Certificate) error {
	var (
		privateKey    interface{}
		privateKeyRaw interface{}
		privateKeyErr error
	)
	expandedPrivateKey, err := homedir.Expand(genAddToSSHAgent)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filename": genAddToSSHAgent,
		}).Error("Could not expand path")
		return fmt.Errorf("Could not expand path: %w", err)
	}

	privatKeyBytes, err := ioutil.ReadFile(expandedPrivateKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filename": expandedPrivateKey,
		}).Error("Could not read private key")
		return fmt.Errorf("Could not read private key: %w", err)
	}

	passPhrase, err := genGetPrivateKeyPassphrase(privatKeyBytes)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filename": expandedPrivateKey,
		}).Error("Could not get passphrase for private key")
		return fmt.Errorf("Could not get passphrase for private key: %w", err)
	}

	privateKeyRaw, privateKeyErr = sshkeys.ParseEncryptedRawPrivateKey(privatKeyBytes, passPhrase)

	if privateKeyErr != nil {
		log.WithFields(log.Fields{
			"error":    privateKeyErr,
			"filename": expandedPrivateKey,
		}).Error("Could not parse private key")
		return fmt.Errorf("Could not parse private key: %w", err)
	}

	privateKey = genCastInterfaceToPrimaryKeyInterface(privateKeyRaw)
	if privateKey == nil {
		log.WithFields(log.Fields{
			"error":    privateKeyErr,
			"filename": expandedPrivateKey,
		}).Error("Unknown private key format")
		return fmt.Errorf("Unknown private key format: %w", err)
	}

	agentAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if agentAuthSock == "" {
		log.WithFields(log.Fields{
			"error":    privateKeyErr,
			"filename": expandedPrivateKey,
		}).Error("SSH_AUTH_SOCK environment variable empty")
		return errors.New("SSH_AUTH_SOCK environment variable empty")
	}
	agentSock, err := net.Dial("unix", agentAuthSock)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    privateKeyErr,
			"filename": expandedPrivateKey,
		}).Error("ssh-agent is not working on SSH_AUTH_SOCK socket")
		return fmt.Errorf("ssh-agent is not working on SSH_AUTH_SOCK socket: %w", err)
	}
	defer agentSock.Close()

	agentKeyring := agent.NewClient(agentSock)

	t := time.Unix(int64(cert.ValidBefore), 0)
	lifetime := t.Sub(time.Now()).Seconds()

	if privateKey == nil {
		log.WithFields(log.Fields{
			"error":    privateKeyErr,
			"filename": expandedPrivateKey,
		}).Error("Unknown private key format")
		return errors.New("Unknown private key format")
	}

	pubcert := agent.AddedKey{
		PrivateKey:   privateKey,
		Certificate:  cert,
		Comment:      fmt.Sprintf("%s [Expires %s]", cert.KeyId, t),
		LifetimeSecs: uint32(lifetime),
	}

	err = agentKeyring.Add(pubcert)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filename": expandedPrivateKey,
		}).Error("Unable to add cert to ssh agent")
		return fmt.Errorf("Unable to add cert to ssh agent: %w", err)
	}

	log.WithFields(log.Fields{
		"cert ID":     cert.KeyId,
		"private key": expandedPrivateKey,
		"valid until": time.Unix(int64(cert.ValidBefore), 0).UTC(),
	}).Info("Cetificate successfully added to ssh-agent")

	return nil
}
