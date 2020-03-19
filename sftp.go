// Helpers to make file management less painful
//
// Not go-routine safe
package sshutils

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
)

func randomName() string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err == nil {
		return base32.StdEncoding.EncodeToString(b[:])
	}
	return ""
}

func (c *Client) getSFTP() (*sftp.Client, error) {
	if c.SFTP != nil {
		return c.SFTP, nil
	}
	s, err := sftp.NewClient(c.Client)
	if err != nil {
		return nil, fmt.Errorf("SFTP client: %w", err)
	}
	c.SFTP = s
	return s, nil
}

func (c *Client) TempPath(base, pfx string) (string, error) {
	s, err := c.getSFTP()
	if err != nil {
		return "", err
	}

	// Create temporary execution dir
	rn := randomName()
	if rn == "" {
		return "", errors.New("Failed to generate random path name")
	}

	s.Mkdir(base)
	rpath := base + "/" + pfx + rn

	// If path already exists, something is really wrong
	if err := s.Mkdir(rpath); err != nil {
		return "", err
	}

	return rpath, nil
}

func (c *Client) TransferFile(local, remote string) error {
	s, err := c.getSFTP()
	if err != nil {
		return err
	}

	l, err := os.Open(local)
	if err != nil {
		return err
	}
	defer l.Close()

	if r, err := s.Create(remote); err != nil {
		return err
	} else {
		defer r.Close()
		_, err := io.Copy(r, l)
		if err != nil {
			return err
		}
	}
	return nil
}
