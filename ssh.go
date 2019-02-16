package sshutils

import (
	"bytes"
	"log"
	"os"
	"sync"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var lock sync.Mutex

var ptyModes = ssh.TerminalModes{
	ssh.ECHO:          0,
	ssh.TTY_OP_ISPEED: 14400,
	ssh.TTY_OP_OSPEED: 14400,
}

// TODO: this will eat the last line if incomplete
type RealTimePrinter struct {
	prefix []byte
	buf    []byte
}

func (p *RealTimePrinter) Write(b []byte) (int, error) {
	p.buf = append(p.buf, b...)
	for {
		c := bytes.Index(p.buf, []byte{'\n'})
		if c < 0 {
			break
		}
		w := append(p.prefix, p.buf[:c+1]...)
		lock.Lock()
		if n, err := os.Stdout.Write(w); err != nil {
			lock.Unlock()
			return n, err
		}
		lock.Unlock()
		p.buf = p.buf[c+1:]
	}
	return len(b), nil
}

func (p *RealTimePrinter) End() {
	// Flush last bytes if data didn't end with a new-line
	if len(p.buf) > 0 {
		p.Write([]byte{'\n'})
	}
}

type Client struct {
	*ssh.Client
	ses  *ssh.Session
	SFTP *sftp.Client

	Name   string
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

type Runner interface {
	Prepare(c *Client, h string, args []string) (string, error)
	Clean(c *Client, h string) error
}

type (
	//RunFunc  func(c *Client, h string, a []string) error
	WaitChan chan struct{}
)

func (c *Client) Run(cmd string) error {
	if c.ses == nil {
		s, err := c.NewSession() // TODO: never closed
		if err != nil {
			return err
		}
		// This has the side-effect of killing the process if the user
		// disconnects; add flag?
		if err := s.RequestPty("xterm", 80, 40, ptyModes); err != nil {
			log.Println(c.Name, "request-pty error:", err)
		}
		r := &RealTimePrinter{[]byte(c.Name + "> "), nil}
		defer r.End()
		s.Stdout = r
		s.Stderr = r
		c.ses = s
	}

	return c.ses.Run(cmd)
}

var (
	auths       []ssh.AuthMethod
	defaultUser string
)

func Run(h InventoryHost, kh *KnownHosts, r Runner, args []string) WaitChan {
	c := make(WaitChan, 1)
	go func() {
		defer close(c)
		u := h.User
		if u == "" {
			u = defaultUser
		}
		hka := kh.GetHKA(h.canonical())
		cfg := &ssh.ClientConfig{
			User:              u,
			Auth:              auths,
			HostKeyAlgorithms: hka,
			HostKeyCallback:   kh.VerifyKey,
		}
		//log.Println("HKA for", h.canonical(), "=", cfg.HostKeyAlgorithms)
		host := h.dialer()
		c, err := ssh.Dial("tcp", host, cfg)
		if err != nil {
			log.Println(h.Name, "failed:", err)
			return
		}
		defer c.Close()

		cli := &Client{Client: c, Name: h.Name}
		defer r.Clean(cli, host)
		cmd, err := r.Prepare(cli, host, args)
		if err != nil {
			log.Println(h.Name, "failed:", err)
		}

		if err := cli.Run(cmd); err != nil {
			log.Println(h.Name, "failed:", err)
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
