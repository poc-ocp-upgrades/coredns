package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func testExternalPluginCompile(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := addExamplePlugin(); err != nil {
		t.Fatal(err)
	}
	defer run(t, gitReset)
	if _, err := run(t, goGet); err != nil {
		t.Fatal(err)
	}
	if _, err := run(t, goGen); err != nil {
		t.Fatal(err)
	}
	if _, err := run(t, goBuild); err != nil {
		t.Fatal(err)
	}
	out, err := run(t, coredns)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "dns.example") {
		t.Fatal("Plugin dns.example should be there")
	}
}
func run(t *testing.T, c *exec.Cmd) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.Dir = ".."
	out, err := c.Output()
	if err != nil {
		return nil, fmt.Errorf("Run: failed to run %s %s: %q", c.Args[0], c.Args[1], err)
	}
	return out, nil
}
func addExamplePlugin() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f, err := os.OpenFile("../plugin.cfg", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(example)
	return err
}

var (
	goBuild		= exec.Command("go", "build")
	goGen		= exec.Command("go", "generate")
	goGet		= exec.Command("go", "get", "github.com/coredns/example")
	gitReset	= exec.Command("git", "checkout", "core/*")
	coredns		= exec.Command("./coredns", "-plugins")
)

const example = "1001:example:github.com/coredns/example"
