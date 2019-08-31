package cmd

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"github.com/le0pard/certonid/utils"
)

func genAddCertToAgent(cert *ssh.Certificate) error {
	var(
		privateKey ssh.Signer
	)
	expandedPrivateKey, err := homedir.Expand(genAddToSSHAgent)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filename": genAddToSSHAgent,
		}).Error("Could not expand path")
		return err
	}

	privatKeyBytes, err := ioutil.ReadFile(expandedPrivateKey)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filename": expandedPrivateKey,
		}).Error("Could not read private key")
		return err
	}

	passphrase, ok := utils.GetENV("PRIVATE_KEY_PASSPHRASE")
	if ok {
		privateKey, err = ssh.ParsePrivateKeyWithPassphrase(privatKeyBytes, []byte(passphrase))
	} else {
		privateKey, err = ssh.ParsePrivateKey(privatKeyBytes)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filename": expandedPrivateKey,
		}).Error("Could not parse private key")
		return err
	}

	agentKeyring := agent.NewKeyring()

	t := time.Unix(int64(cert.ValidBefore), 0)
	lifetime := t.Sub(time.Now()).Seconds()

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
		panic(reflect.TypeOf(privateKey))
		return err
	}

	log.WithFields(log.Fields{
		"error":       err,
		"filename":    expandedPrivateKey,
		"valid until": time.Unix(int64(cert.ValidBefore), 0).UTC(),
	}).Info("Cetificate successfully added to ssh-agent")

	return nil
}
