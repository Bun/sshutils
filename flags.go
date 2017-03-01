package sshutils

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"io/ioutil"
	"net"
	"os"
	"os/user"
	"strings"
)

func loadAgent() {
	if a, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		am := ssh.PublicKeysCallback(agent.NewClient(a).Signers)
		auths = append(auths, am)
	}
}

func filteredInventory(all []InventoryHost, hosts string) []InventoryHost {
	if hosts == "all" || hosts == "*" {
		return all
	}

	var inv []InventoryHost
	wcs := wildcards(strings.Split(hosts, ";"))

	for _, h := range all {
		m := false
		for _, wc := range wcs {
			if wc.MatchString(h.Name) {
				m = true
				break
			}
		}
		if m {
			inv = append(inv, h)
		}
	}

	return inv
}

// TODO: -i inv
// TODO: -u user
func ParseFlags() ([]InventoryHost, []string) {
	loadAgent()

	u, _ := user.Current()
	defaultUser = u.Username

	b, err := ioutil.ReadFile("inv")
	if err != nil {
		panic(err)
	}

	inv := parseInventory(b)
	if len(os.Args) > 1 {
		// Extract relevant hosts
		inv = filteredInventory(inv, os.Args[1])
	}

	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	return inv, args
}
