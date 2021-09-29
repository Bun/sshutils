package sshutils

import (
	"bytes"
	"testing"
)

func TestSSHConfigParser(t *testing.T) {
	sc := parseSSHConfig(bytes.NewBuffer([]byte(`
host foo
 hostname test
host bar
 hostname tset`)))

	t.Logf("%#v", sc)
	if sc.HostAlias["foo"] != "test" {
		t.Error()
	}
	if sc.HostAlias["bar"] != "tset" {
		t.Error()
	}

	// Should not crash
	t.Logf("%#v", LoadSSHConfig())
}
