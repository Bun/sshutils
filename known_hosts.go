package sshutils

// TODO: optimize some lookups based on hostname

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"

	"golang.org/x/crypto/ssh"
)

type KnownHostEntry struct {
	E         string
	Key, Type string
	MatchHost func(e, name string) bool
}

type KnownHosts struct {
	Entries []KnownHostEntry
}

// Load known_hosts entries
func LoadKnownHosts() (kh *KnownHosts) {
	fn := path.Join(os.Getenv("HOME"), ".ssh/known_hosts")
	kh = &KnownHosts{}
	if d, err := ioutil.ReadFile(fn); err == nil {
		kh.Parse(string(d))
	} else {
		panic(err) // TODO
	}
	return kh
}

func (kh *KnownHosts) Parse(s string) {
	// host, key-type, key
	for _, line := range strings.Split(s, "\n") {
		parts := strings.Split(line, " ")
		if len(parts) < 3 {
			continue
		}
		pk := parts[1] + " " + parts[2]
		if strings.HasPrefix(parts[0], "|1|") {
			kh.Entries = append(kh.Entries, KnownHostEntry{
				parts[0], pk, parts[1], HashedHMACSHA1Entry})
		} else {
			kh.Entries = append(kh.Entries, KnownHostEntry{
				parts[0], pk, parts[1], PlainHostsEntry})
		}
	}
}

func (kh *KnownHosts) VerifyKey(hostname string, remote net.Addr, key ssh.PublicKey) error {
	hostname = canonicalHost(hostname)
	actual := pubkey(key)
	offered := key.Type()
	var seen []string
	for _, entry := range kh.Entries {
		if entry.MatchHost(entry.E, hostname) {
			seen = append(seen, entry.Type)
			if entry.Type != offered {
				continue
			} else if entry.Key == actual {
				return nil
			}
			return fmt.Errorf("%s: Unexpected host key: %s", hostname, actual)
		}
	}
	if len(seen) > 0 {
		return fmt.Errorf("%s: Offered %s, know only %s", hostname, offered, seen)
	}
	return fmt.Errorf("%s: Hostname not in known_hosts", hostname)
}

// TODO: prefer order
func (kh *KnownHosts) GetHKA(hostname string) (hka []string) {
	for _, entry := range kh.Entries {
		if entry.MatchHost(entry.E, hostname) {
			hka = append(hka, entry.Type)
		}
	}
	return
	//hka = []string{
	//	ssh.KeyAlgoED25519,
	//	ssh.KeyAlgoECDSA256,
	//	ssh.KeyAlgoECDSA384,
	//	ssh.KeyAlgoRSA,
	//}
}

func PlainHostsEntry(cached, name string) bool {
	for _, n := range strings.Split(cached, ",") {
		if n == name {
			return true
		}
	}
	return false
}

// HMAC-SHA1: |1|salt|digest|
// TODO: check how OpenSSH deals with collisions
func HashedHMACSHA1Entry(e, name string) bool {
	parts := strings.Split(e, "|")
	if len(parts) < 4 || parts[0] != "" || parts[1] != "1" {
		fmt.Printf("invalid entry %q\n", parts)
		return false
	}
	salt, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		println("invalid salt")
		return false
	}
	actual, err := base64.StdEncoding.DecodeString(parts[3])
	if err != nil {
		println("invalid digest")
		return false
	}
	h := hmac.New(sha1.New, salt)
	h.Write([]byte(name))
	expected := h.Sum(nil)
	return hmac.Equal(expected, actual)
}
