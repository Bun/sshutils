package sshutils

import (
	"flag"
	"fmt"
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
	flagNokey = flag.Bool("nk", false, "Run even if no keys are loaded in agent")
	flagInv   = flag.String("i", "hosts.yaml", "Inventory filename")
	flagUser  = flag.String("u", "", "Default username")
)

func ParseFlags() ([]Target, []string) {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage:", os.Args[0], "[flags] <host-pattern> <command>")
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if err := loadAgent(*flagNokey); err != nil {
		log.Fatalln("ssh-agent error:", err)
	}

	if *flagUser != "" {
		forceUser = *flagUser
	}
	u, _ := user.Current()
	defaultUser = u.Username

	inv, err := parseInventory(*flagInv)
	if err != nil {
		log.Fatalln("Inventory:", err)
	}

	return filteredInventory(inv, args[0]), args[1:]
}
