package gocron_server_test

import (
	"testing"

	gocron_server "github.com/strnophix/gocron-server/pkg"
)

func TestUnitExecutableCommand(t *testing.T) {
	cmd := gocron_server.NewUnitExecCmd("echo hi")
	out, err := cmd.Call()
	if err != nil {
		t.Fatalf("Should not have errored on a simple echo")
	}

	if out != "hi\n" {
		t.Fatalf("Execution should have returned \"hi\" but gave: %s", out)
	}

	c := Counter{Current: 1}
	fn := gocron_server.NewUnitExecFn(c.Increment)
	out, _ = fn.Call()
	if out != "2" {
		t.Fatalf("Execution of unit should have incremented counter but got: %d", c.Current)
	}

	cmd = gocron_server.NewUnitExecCmd("false")
	_, err = cmd.Call()
	if err == nil {
		t.Fatalf("Should have returned an error on a non-zero exit")
	}
}
