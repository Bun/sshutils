package main

import (
	"awoo.nl/sshutils"

	"strings"
)

type Run struct{}

func (_ *Run) Prepare(c *sshutils.Client, h string, args []string) (string, error) {
	// HEADSUP: Commands are expected to be escaped by the user if necessary
	return strings.Join(args, " "), nil
}

func (_ *Run) Clean(c *sshutils.Client, h string) error {
	return nil
}

func main() {
	var ws []sshutils.WaitChan
	hosts, args := sshutils.ParseFlags()
	kh := sshutils.LoadKnownHosts()
	sc := sshutils.LoadSSHConfig()

	for _, h := range hosts {
		if h.Host == "" {
			h.Host = sc.HostAlias[h.Name]
		}
		ws = append(ws, sshutils.Run(h, kh, &Run{}, args))
	}

	sshutils.WaitAll(ws)
}
