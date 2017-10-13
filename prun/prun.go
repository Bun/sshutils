package main

import (
	"path"
	"strings"

	"awoo.nl/sshutils"
)

type Run struct {
	rpath, rfn string
}

func (r *Run) Prepare(c *sshutils.Client, h string, args []string) (string, error) {
	var err error
	fn := args[0]
	r.rpath, err = c.TempPath(".cache", "prun-")
	if err != nil {
		return "", err
	}
	r.rfn = path.Join(r.rpath, path.Base(fn))
	if err := c.TransferFile(fn, r.rfn); err != nil {
		return "", err
	}
	if err := c.SFTP.Chmod(r.rfn, 0755); err != nil {
		return "", err
	}

	// TODO: escape rfn!
	cmd := r.rfn
	if len(args) > 1 {
		cmd += " " + strings.Join(args[1:], " ")
	}
	return cmd, nil
}

func (r *Run) Clean(c *sshutils.Client, h string) error {
	if c.SFTP != nil {
		c.SFTP.Remove(r.rfn)
		c.SFTP.RemoveDirectory(r.rpath)
	}
	return nil
}

func main() {
	var ws []sshutils.WaitChan
	hosts, args := sshutils.ParseFlags()
	kh := sshutils.LoadKnownHosts()

	for _, h := range hosts {
		ws = append(ws, sshutils.Run(h, kh, &Run{}, args))
	}

	sshutils.WaitAll(ws)
}
