package main

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"github.com/Bun/sshutils"
	"github.com/pkg/sftp"
	"io"
	"os"
	"path"
	"strings"
)

func randomName() string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err == nil {
		return base32.StdEncoding.EncodeToString(b[:])
	}
	return ""
}

// TODO: audit
// TODO: test error handling
// TODO: resource leakage
func run(c *sshutils.Client, h string, args []string) error {
	s, err := sftp.NewClient(c.Client)
	if err != nil {
		return err
	}
	defer s.Close()

	// Create temporary execution dir
	rn := randomName()
	if rn == "" {
		return errors.New("Failed to generate random path name")
	}
	// TODO: variable path
	base := ".cache"
	s.Mkdir(base)
	rpath := base + "/prun-" + rn
	// If path already exists, something is really wrong
	if err := s.Mkdir(rpath); err != nil {
		return err
	}

	// Transfer file
	fn := args[0]
	l, err := os.Open(fn)
	if err != nil {
		return err
	}

	b := path.Base(fn)
	rfn := rpath + "/" + b
	if r, err := s.Create(rfn); err != nil {
		return err
	} else {
		_, err := io.Copy(r, l)
		if err != nil {
			return err
		}
		r.Close()
	}
	l.Close()
	if err := s.Chmod(rfn, 0755); err != nil {
		return err
	}

	// TODO: escape rfn!
	cmd := rfn
	if len(args) > 1 {
		cmd += " " + strings.Join(args[1:], " ")
	}
	err = c.Run(cmd)
	c.Output()

	// Erase temporary paths
	s.Remove(rfn)
	s.RemoveDirectory(rpath)

	return err
}

func main() {
	var ws []sshutils.WaitChan
	hosts, args := sshutils.ParseFlags()

	for _, h := range hosts {
		ws = append(ws, sshutils.Run(h, run, args))
	}

	sshutils.WaitAll(ws)
}
