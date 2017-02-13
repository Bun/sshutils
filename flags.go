package sshutils

import (
	"golang.org/x/crypto/ssh"

	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

func loadKey(fn string) (ssh.Signer, error) {
	key, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(key)
}

// TODO
func loadKeys() {
	u, _ := user.Current()
	// XXX
	defaultUser = u.Username
	d := u.HomeDir + "/.ssh/"
	fis, err := ioutil.ReadDir(d)
	if err != nil {
		panic(err)
	}
	for _, fi := range fis {
		fn := fi.Name()
		if !strings.HasSuffix(fn, ".pub") {
			if s, err := loadKey(d + fn); err == nil {
				auths = append(auths, ssh.PublicKeys(s))
			}
		}
	}
}

func match(name, wc string) bool {
	// TODO
	return name == wc
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

func ParseFlags() ([]InventoryHost, []string) {
	loadKeys()

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
