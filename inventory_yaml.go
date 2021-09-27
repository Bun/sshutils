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

type yamlInventory struct {
	All struct {
		Hosts map[string]struct {
			Host    string `yaml:"ansible_host"`
			User    string `yaml:"ansible_user"`
			SSHPort string `yaml:"ansible_ssh_port"`
			SSHHost string `yaml:"ansible_ssh_host"`
		} `yaml:"hosts"`
		Vars struct {
			User string `yaml:"ansible_user"`
		} `yaml:"vars"`
	} ` yaml:"all"`
}

func parseInventoryYaml(fname string) (inv Inventory, err error) {
	var req yamlInventory
	if err = readFile(fname, &req); err != nil {
		return Inventory{}, err
	}
	for name, host := range req.All.Hosts {
		h := Target{
			Name: name,
			User: picks(host.User, req.All.Vars.User),
			Host: picks(host.SSHHost, host.Host, name),
			Port: host.SSHPort,
		}
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
