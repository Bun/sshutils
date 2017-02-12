package sshutils

import (
	"golang.org/x/crypto/ssh"

	"bytes"
	"fmt"
	"strings"
)

type Client struct {
	*ssh.Client
	ses *ssh.Session

	Name   string
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

type (
	RunFunc  func(c *Client, h string, a []string) error
	WaitChan chan struct{}
)

func (c *Client) Run(cmd string) error {
	if c.ses == nil {
		s, err := c.NewSession() // TODO: never closed
		if err != nil {
			return err
		}
		s.Stdout = &c.Stdout
		s.Stderr = &c.Stderr
		c.ses = s
	}

	return c.ses.Run(cmd)
}

func (c *Client) Output() {
	fmt.Printf("%s |\n", c.Name)

	o := strings.TrimRight(c.Stdout.String(), " \t\r\n")
	if o != "" {
		fmt.Printf("O: %s\n", o)
	}
	o = strings.TrimRight(c.Stderr.String(), " \t\r\n")
	if o != "" {
		fmt.Printf("E: %s\n", o)
	}
}

var (
	auths       []ssh.AuthMethod
	defaultUser string
	hka         = []string{
		ssh.KeyAlgoED25519,
		ssh.KeyAlgoECDSA384,
		ssh.KeyAlgoECDSA256,
		ssh.KeyAlgoRSA,
	}
)

func Run(h InventoryHost, rf RunFunc, args []string) WaitChan {
	c := make(WaitChan, 1)
	go func() {
		defer close(c)
		host := h.Host
		u := h.User
		if host == "" {
			host = h.Name
		}
		if h.Port != "" {
			host += ":" + h.Port
		} else {
			host += ":22"
		}
		if u == "" {
			u = defaultUser
		}
		cfg := &ssh.ClientConfig{
			User:              u,
			Auth:              auths,
			HostKeyAlgorithms: hka,
		}
		c, err := ssh.Dial("tcp", host, cfg)
		if err != nil {
			panic(err)
		}
		defer c.Close()
		err = rf(&Client{Client: c, Name: h.Name}, host, args)
		if err != nil {
			fmt.Println(h.Name, "failed:", err)
		}
	}()
	return c
}

func WaitAll(ws []WaitChan) {
	for _, w := range ws {
		ok := true
		for ok {
			select {
			case _, ok = <-w:
			}
		}
	}
}
