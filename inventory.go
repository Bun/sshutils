package sshutils

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Inventory struct {
	Targets []Target
	Groups  map[string][]string
}

type Target struct {
	Name string
	Host string
	Port string
	User string
}

// For comparing hostnames
func (h Target) canonical() string {
	host := h.Host
	if host == "" {
		host = h.Name
	}
	if h.Port != "" && h.Port != "22" {
		return "[" + host + "]:" + h.Port
	}
	if strings.Index(host, ":") >= 0 {
		return "[" + host + "]"
	}
	return host
}

// For dialing
func (h Target) dialer() string {
	host := h.Host
	if host == "" {
		host = h.Name
	}
	if h.Port != "" {
		host += ":" + h.Port
	} else {
		host += ":22"
	}
	return host
}

func prefix(v, s string) (string, bool) {
	if strings.HasPrefix(v, s) {
		return v[len(s):], true
	}
	return "", false
}

func parseInventory(name string) (Inventory, error) {
	if strings.HasSuffix(name, ".yaml") {
		return parseInventoryYaml(name)
	}
	info, rerr := os.Stat(name)
	if rerr == nil && !info.IsDir() {
		return parseInventoryINI(name)
	}
	hy := filepath.Join(name, "hosts.yaml")
	if inv, err := parseInventoryYaml(hy); err == nil || !os.IsNotExist(err) {
		return inv, err
	}
	hi := filepath.Join(name, "hosts")
	if inv, err := parseInventoryINI(hi); err == nil || !os.IsNotExist(err) {
		return inv, err
	}
	return Inventory{}, rerr
}

// TODO: support groups
func parseInventoryINI(fname string) (inv Inventory, err error) {
	buf, err := os.ReadFile(fname)
	if err != nil {
		return
	}
	// TODO: support groups
	lines := bytes.Split(buf, []byte("\n"))
	for _, bline := range lines {
		line := string(bline)
		if c := strings.Index(line, "#"); c >= 0 {
			line = line[:c]
		}
		if strings.Trim(line, " \t") == "" {
			continue
		}
		parts := strings.Split(line, " ")
		h := Target{Name: parts[0]}
		for _, p := range parts[1:] {
			if v, ok := prefix(p, "ansible_host="); ok {
				h.Host = v
			} else if v, ok := prefix(p, "ansible_ssh_host="); ok {
				h.Host = v
			} else if v, ok := prefix(p, "ansible_ssh_port="); ok {
				h.Port = v
			} else if v, ok := prefix(p, "ansible_ssh_user="); ok {
				h.User = v
			} else if v, ok := prefix(p, "ansible_user="); ok {
				h.User = v
			}
		}
		inv.Targets = append(inv.Targets, h)
	}
	return
}
