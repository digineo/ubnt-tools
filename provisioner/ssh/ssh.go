package ssh

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	"encoding/pem"

	"strings"

	"github.com/digineo/goldflags"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type sshSessionCallback func(*ssh.Session) error

// Agent tries to connect with the ssh-agent
func Agent() ssh.AuthMethod {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		log.Printf("[ssh.Agent] SSH_AUTH_SOCK is not defined or empty")
		return nil
	}

	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err == nil {
		log.Printf("[ssh.Agent] Couldn't connect to SSH agent: %v", err)
		return nil
	}

	agent := agent.NewClient(sshAgent)
	keys, err := agent.List()
	if err != nil {
		log.Printf("[ssh.Agent] Listing keys error'ed: %v", err)
	} else {
		log.Printf("Keys: %v", keys)
	}
	return ssh.PublicKeysCallback(agent.Signers)
}

// ReadPrivateKey tries to read an SSH private key file.
func ReadPrivateKey(keyPath, password string) (auth ssh.AuthMethod, ok bool) {
	keyFile, err := goldflags.ExpandPath(keyPath)
	if err != nil {
		log.Printf("[ssh.ReadPrivateKey] Could not expand %s: %v", keyPath, err)
		return
	}

	if !goldflags.PathExist(keyFile) {
		log.Printf("[ssh.ReadPrivateKey] Keyfile %s not found", keyFile)
		return
	}

	keyPEM, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Printf("[ssh.ReadPrivateKey] Could not read %s: %v", keyFile, err)
		return
	}

	block, _ := pem.Decode(keyPEM)
	if block == nil {
		log.Printf("[ssh.ReadPrivateKey] No key found in %s", keyFile)
		return
	}

	keyFrom := block.Bytes
	if strings.Contains(block.Headers["Proc-Type"], "ENCRYPTED") {
		keyFrom, err = x509.DecryptPEMBlock(block, []byte(password))
		if err != nil {
			log.Printf("[ssh.ReadPrivateKey] Error decrypting %s: %v", keyFile, err)
			return
		}
	}

	key, err := getKey(block.Type, keyFrom)
	if err != nil {
		log.Printf("[ssh.ReadPrivateKey] %s: %v", keyFile, err)
		return
	}

	sign, err := ssh.NewSignerFromKey(key)
	if err != nil {
		log.Printf("[ssh.ReadPrivateKey] %s: %v", keyFile, err)
		return
	}
	return ssh.PublicKeys(sign), true
}

func getKey(typ string, b []byte) (interface{}, error) {
	switch typ {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(b)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(b)
	case "DSA PRIVATE KEY":
		return ssh.ParseDSAPrivateKey(b)
	default:
		return nil, fmt.Errorf("unsupported key type %q", typ)
	}
}

// WithinSession executes a callback function within a new SSH session of
// the given client.
func WithinSession(client *ssh.Client, callback sshSessionCallback) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return callback(session)
}

// UploadFile uploads a local file to the remote. Please avoid funky
// remote file names, since there's no protection against command injection
func UploadFile(client *ssh.Client, localName string, remoteName string) error {
	return WithinSession(client, func(s *ssh.Session) error {
		writer, err := s.StdinPipe()
		if err != nil {
			return err
		}
		defer writer.Close()

		buf, err := ioutil.ReadFile(localName)
		if err != nil {
			return err
		}
		log.Printf("[ssh.UploadFile local-file] %s, %d bytes", localName, len(buf))

		rdir := filepath.Dir(remoteName)
		log.Printf("[ssh.UploadFile remote-dir] %s", rdir)

		rfile := filepath.Base(remoteName)
		log.Printf("[ssh.UploadFile remote-file] %s", rfile)

		var so, se bytes.Buffer
		s.Stdout = &so
		s.Stderr = &se

		cmd := fmt.Sprintf("/usr/bin/scp -t %s", rdir) // danger!
		log.Printf("[ssh.UploadFile command] %s", cmd)

		if err := s.Start(cmd); err != nil {
			return err
		}

		content := string(buf)
		log.Printf("[ssh.UploadFile uploading] %d bytes", len(content))

		// https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works
		fmt.Fprintln(writer, "C0644", len(content), rfile)
		fmt.Fprint(writer, content)
		fmt.Fprint(writer, "\x00")
		writer.Close()

		if err := s.Wait(); err != nil {
			log.Println("[ssh.UploadFile] waiting failed")
			log.Printf("[ssh.UploadFile stderr] %s", se.String())
			log.Printf("[ssh.UploadFile stdout] %s", so.String())
			return err
		}

		return nil
	})
}

// ExecuteCommand executes a command in a new SSH session.
func ExecuteCommand(client *ssh.Client, cmd string) (string, error) {
	var output string
	sessionErr := WithinSession(client, func(s *ssh.Session) error {
		var so, se bytes.Buffer
		s.Stdout = &so
		s.Stderr = &se

		if err := s.Run(cmd); err != nil {
			log.Printf("[executeCommand] %s failed", cmd)
			log.Printf("[executeCommand stderr] %s", se.String())
			log.Printf("[executeCommand stdout] %s", so.String())
			return err
		}

		output = se.String()
		return nil
	})

	return output, sessionErr
}
