package sshutils

import (
	"encoding/base64"
	"log"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
)

var wildcardExpr = strings.NewReplacer(
	"\\*", ".*",
	"\\?", ".?",
)

// TODO: error reporting
func wildcards(wcs []string) (rs []*regexp.Regexp) {
	for _, wc := range wcs {
		wc := wildcardExpr.Replace(regexp.QuoteMeta(wc))
		if r, err := regexp.Compile("^" + wc + "$"); err == nil {
			rs = append(rs, r)
		} else {
			log.Println("Wildcard failed:", err)
		}
	}
	return
}

// Format a host to be in `[host]:port` or `host` format
func canonicalHost(s string) string {
	var host, port string
	if strings.HasPrefix(s, "[") {
		c := strings.Index(s, "]")
		if c < 0 {
			return s
		}
		host = s[1:c]
		if strings.HasPrefix(s[c+1:], ":") {
			// Junk at the end shouldn't result in a valid match
			port = s[c+2:]
		}
	} else if c := strings.LastIndex(s, ":"); c > 0 {
		host = s[:c]
		port = s[c+1:]
	} else {
		host = s
	}
	if port == "22" {
		port = ""
	}
	if port != "" {
		return "[" + host + "]:" + port
	}
	return host
}

func pubkey(pk ssh.PublicKey) string {
	t := pk.Type()
	b := base64.StdEncoding.EncodeToString(pk.Marshal())
	return t + " " + b
}
