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

func ParseFlags() ([]InventoryHost, []string) {
	loadKeys()

	b, err := ioutil.ReadFile("inv")
	if err != nil {
		panic(err)
	}

	inv := parseInventory(b)
	// <extract relevant hosts>
	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	return inv, args
}
