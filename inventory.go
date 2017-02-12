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
		parts := strings.Split(line, " ")
		if parts[0] == "" {
			continue
		}
		h := InventoryHost{Name: parts[0]}
		for _, p := range parts[1:] {
			if v, ok := prefix(p, "ansible_ssh_host="); ok {
				h.Host = v
			} else if v, ok := prefix(p, "ansible_ssh_port="); ok {
				h.Port = v
			} else if v, ok := prefix(p, "ansible_ssh_user="); ok {
				h.User = v
			}
		}
		inv = append(inv, h)
	}
	return inv
}
