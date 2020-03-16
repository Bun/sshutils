package sshutils

import (
	"gopkg.in/yaml.v3"
)

// We support a limited subset of Ansible's YAML config
type (
	YInventoryHost map[string]interface{}
	YInventoryVars map[string]interface{}

	// An InventoryChild node is used to bundle vars and hosts together
	YInventoryChild struct {
		Hosts map[string]YInventoryHost `yaml:"hosts"`
		Vars  YInventoryVars            `yaml:"vars"`
	}

	YInventory struct {
		Children map[string]YInventoryChild `yaml:"children"`
		Hosts    map[string]YInventoryHost  `yaml:"hosts"`
		Vars     YInventoryVars             `yaml:"vars"`
	}

	InventoryYaml struct {
		All YInventory `yaml:"all"`
	}
)

func yamlInv(b []byte) (Inventory, error) {
	var iyaml InventoryYaml
	err := yaml.Unmarshal(b, &iyaml)
	if err != nil {
		return Inventory{}, err
	}

	ansibleUser, _ := iyaml.All.Vars["ansible_user"].(string)
	ansiblePort, _ := iyaml.All.Vars["ansible_port"].(string)

	groups := map[string][]string{}
	hosts := map[string]Target{}

	for group, data := range iyaml.All.Children {
		u, _ := data.Vars["ansible_user"].(string)
		p, _ := data.Vars["ansible_port"].(string)
		for h := range data.Hosts {
			groups[group] = append(groups[group], h)
			hosts[h] = Target{
				Name: h,
				User: u,
				Port: p,
			}
		}
	}

	for host, data := range iyaml.All.Hosts {
		t := hosts[host]
		t.Name = host

		u := t.User
		if u == "" {
			u = ansibleUser
		}
		if hu, ok := data["ansible_user"].(string); ok {
			u = hu
		}
		if u != "" {
			t.User = u
		}

		h, _ := data["ansible_host"].(string)
		t.Host = h

		p, _ := data["ansible_port"].(string)
		if p == "" {
			p = ansiblePort
		}
		if p != "" {
			t.Port = p
		}

		hosts[host] = t
	}

	var targets []Target
	for _, t := range hosts {
		targets = append(targets, t)
	}

	return Inventory{
		Targets: targets,
		Groups:  groups,
	}, nil
}
