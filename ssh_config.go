package sshutils

// TODO: optimize some lookups based on hostname

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path"
	"strings"
)

type SSHConfig struct {
	HostAlias map[string]string
}

func LoadSSHConfig() SSHConfig {
	fn := path.Join(os.Getenv("HOME"), ".ssh/config")
	f, err := os.Open(fn)
	if err == nil {
		defer f.Close()
		return parseSSHConfig(f)
	}
	return SSHConfig{}
}

func parseSSHConfig(r io.Reader) (sc SSHConfig) {
	sc.HostAlias = map[string]string{}
	br := bufio.NewReader(r)
	var hosts []string
	for {
		line, _, err := br.ReadLine()
		line = bytes.ToLower(bytes.Trim(line, " \t"))
		if bytes.HasPrefix(line, []byte("host ")) {
			hosts = strings.Split(string(line[5:]), " ")
		}
		if bytes.HasPrefix(line, []byte("hostname ")) {
			hostname := strings.Trim(string(line[9:]), " \t")
			for _, h := range hosts {
				sc.HostAlias[h] = hostname
			}
		}
		if err == io.EOF {
			break
		}
	}
	return
}
