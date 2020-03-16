package sshutils

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func loadAgent(keysOptional bool) error {
	a, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return err
	}
	agent := agent.NewClient(a)
	keys, err := agent.List()
	if err != nil {
		return err
	} else if !keysOptional && len(keys) == 0 {
		// Could be some edge case where this isn't an issue
		return fmt.Errorf("Agent has no keys loaded")
	}

	am := ssh.PublicKeysCallback(agent.Signers)
	auths = append(auths, am)
	return nil
}

func filteredInventory(all Inventory, hosts string) []Target {
	if hosts == "all" || hosts == "*" {
		return all.Targets
	}

	// TODO: support groups

	var inv []Target
	wcs := wildcards(strings.Split(hosts, ";"))

	for _, h := range all.Targets {
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

var (
	flagNokey = flag.Bool("nk", false, "Keys are not required")
	flagInv   = flag.String("i", "hosts.yaml", "Inventory filename")
	flagUser  = flag.String("u", "", "Default username")
)

func ParseFlags() ([]Target, []string) {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatalln("Usage: X host-selector command [args]")
	}

	if err := loadAgent(*flagNokey); err != nil {
		log.Fatalln("ssh-agent error:", err)
	}

	if *flagUser != "" {
		defaultUser = *flagUser
	} else {
		u, _ := user.Current()
		defaultUser = u.Username
	}

	i, err := inv(*flagInv)
	if err != nil {
		log.Fatalln("Inventory:", err)
	}

	return filteredInventory(i, args[0]), args[1:]
}

func inv(fname string) (Inventory, error) {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal("Error reading ", *flagInv, ": ", err)
	}

	if strings.HasSuffix(fname, ".yaml") || strings.HasSuffix(fname, ".yml") {
		return yamlInv(b)
	}

	// Extract relevant hosts
	return parseInventory(b)
}
