package sshutils

import "strings"

type (
	InventoryHost struct {
		Name string
		Host string
		Port string
		User string
	}
)

// For comparing hostnames
func (h InventoryHost) canonical() string {
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
func (h InventoryHost) dialer() string {
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

// TODO: support groups
func parseInventory(buf []byte) (inv []InventoryHost) {
	lines := strings.Split(string(buf), "\n")
	for _, line := range lines {
		if c := strings.Index(line, "#"); c >= 0 {
			line = line[:c]
		}
		if strings.Trim(line, " \t") == "" {
			continue
		}
		parts := strings.Split(line, " ")
		h := InventoryHost{Name: parts[0]}
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
		inv = append(inv, h)
	}
	return inv
}
