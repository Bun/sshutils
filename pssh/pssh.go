package main

import (
	"awoo.nl/sshutils"

	"strings"
)

func run(c *sshutils.Client, h string, args []string) error {
	err := c.Run(strings.Join(args, " "))
	c.Output()
	return err
}

// HEADSUP: Commands are expected to be escaped by the user if necessary
func main() {
	var ws []sshutils.WaitChan
	hosts, args := sshutils.ParseFlags()

	for _, h := range hosts {
		ws = append(ws, sshutils.Run(h, run, args))
	}

	sshutils.WaitAll(ws)
}
