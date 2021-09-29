package sshutils

import (
	"os"

	"gopkg.in/yaml.v3"
)

// readFile reads YAML data from a filepath.
func readFile(fname string, v interface{}) error {
	buf, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, v)
}

type yamlInventory map[string]struct {
	Hosts map[string]struct {
		Host    string `yaml:"ansible_host"`
		User    string `yaml:"ansible_user"`
		SSHPort string `yaml:"ansible_ssh_port"`
		SSHHost string `yaml:"ansible_ssh_host"`
	} `yaml:"hosts"`
	Children map[string]struct{} `yaml:"children"`
	Vars     struct {
		User string `yaml:"ansible_user"`
	} `yaml:"vars"`
}

func parseInventoryYaml(fname string) (inv Inventory, err error) {
	var req yamlInventory
	if err = readFile(fname, &req); err != nil {
		return Inventory{}, err
	}
	inv.Groups = map[string][]string{}
	addToGroup := func(g, h string) {
		for _, v := range inv.Groups[g] {
			if v == h {
				return
			}
		}
		inv.Groups[g] = append(inv.Groups[g], h)
	}
	fallbackUser := req["all"].Vars.User
	hosts := map[string]Target{}

	for group, g := range req {
		for name, host := range g.Hosts {
			addToGroup(group, name)
			h := hosts[name]
			h.Name = picks(h.Name, name)
			h.User = picks(h.User, host.User, g.Vars.User, fallbackUser)
			h.Host = picks(h.Host, host.SSHHost, host.Host)
			h.Port = picks(h.Port, host.SSHPort)
			hosts[name] = h
		}
		//for name, _ := range g.Children {
		// XXX: this makes `name` a subgroup of the current group
		//}
	}
	for _, h := range hosts {
		inv.Targets = append(inv.Targets, h)
	}
	return
}

func picks(args ...string) string {
	for _, arg := range args {
		if arg != "" {
			return arg
		}
	}
	return ""
}
