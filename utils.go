package sshutils

import (
	"golang.org/x/crypto/ssh"

	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"regexp"
	"strings"
)

var wildcardExpr = strings.NewReplacer(
	"\\*", ".*",
	"\\?", ".?",
)

// TODO: error reporting
func wildcards(wcs []string) (rs []*regexp.Regexp) {
	for _, wc := range wcs {
		wc := wildcardExpr.Replace(regexp.QuoteMeta(wc))
		if r, err := regexp.Compile(wc); err == nil {
			rs = append(rs, r)
		} else {
			log.Println("Wildcard failed:", err)
		}
	}
	return
}

func fixHost(s string) string {
	// XXX either [host]:port or host
	var host, port string
	if strings.HasPrefix(s, "[") {
		c := strings.Index(s, "]")
		if c < 0 {
			return s
		}
		host = s[1:c]
		if strings.HasPrefix(s[c+1:], ":") {
			port = s[c+2:]
		}
	} else if c := strings.LastIndex(s, ":"); c > 0 {
		host = s[:c]
		port = s[c+1:]
	} else {
		host = s
	}
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	if port == "" {
		port = "22"
	}
	return host + ":" + port
}

type KnownHosts map[string]string

func (kh KnownHosts) VerifyKey(hostname string, remote net.Addr, key ssh.PublicKey) error {
	expect, ok := kh[hostname]
	if !ok {
		return fmt.Errorf("Unknown host: %s", hostname)
	}
	if actual := pubkey(key); actual != expect {
		return fmt.Errorf("%s: Unexpected host key: %s", hostname, actual)
	}
	return nil
}

func (kh KnownHosts) Parse(s string) {
	for _, line := range strings.Split(s, "\n") {
		parts := strings.Split(line, " ")
		if len(parts) < 3 { //|| !strings.HasPrefix(parts[2], "ssh-") {
			continue
		}
		pk := parts[1] + " " + parts[2]
		for _, host := range strings.Split(parts[0], ",") {
			kh[fixHost(host)] = pk
		}
	}
}

func LoadKnownHosts() (kh KnownHosts) {
	fn := path.Join(os.Getenv("HOME"), ".ssh/known_hosts")
	kh = make(KnownHosts)
	if d, err := ioutil.ReadFile(fn); err == nil {
		kh.Parse(string(d))
	} else {
		panic(err)
	}
	return kh
}

func pubkey(pk ssh.PublicKey) string {
	t := pk.Type()
	b := base64.StdEncoding.EncodeToString(pk.Marshal())
	return t + " " + b
}
