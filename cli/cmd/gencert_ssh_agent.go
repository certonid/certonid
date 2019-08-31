package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/le0pard/certonid/utils"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

func genAddCertToAgent(cert *ssh.Certificate) error {
	var (
		privateKey    ssh.Signer
		privateKeyErr error
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

	_, ok := utils.GetENV("PRIVATE_KEY_PASSPHRASE")
	if ok {
		fmt.Print("SSH Key Passphrase [none]: ")
		passPhrase, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		fmt.Print("\n")
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Not provided passphrase for provate key")
			return err
		}

		privateKey, privateKeyErr = ssh.ParsePrivateKeyWithPassphrase(privatKeyBytes, []byte(passPhrase))
	} else {
		privateKey, privateKeyErr = ssh.ParsePrivateKey(privatKeyBytes)
	}

	if privateKeyErr != nil {
		log.WithFields(log.Fields{
			"error":    privateKeyErr,
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
		return err
	}

	log.WithFields(log.Fields{
		"error":       err,
		"filename":    expandedPrivateKey,
		"valid until": time.Unix(int64(cert.ValidBefore), 0).UTC(),
	}).Info("Cetificate successfully added to ssh-agent")

	return nil
}
