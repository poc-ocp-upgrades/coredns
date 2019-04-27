package test

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/mholt/caddy"
)

var contents = map[string]string{"Kexample.org.+013+45330.key": examplePub, "Kexample.org.+013+45330.private": examplePriv, "example.org.signed": exampleOrg}

const (
	examplePub	= `example.org. IN DNSKEY 256 3 13 eNMYFZYb6e0oJOV47IPo5f/UHy7wY9aBebotvcKakIYLyyGscBmXJQhbKLt/LhrMNDE2Q96hQnI5PdTBeOLzhQ==
`
	examplePriv	= `Private-key-format: v1.3
Algorithm: 13 (ECDSAP256SHA256)
PrivateKey: f03VplaIEA+KHI9uizlemUSbUJH86hPBPjmcUninPoM=
`
)

func TestReadme(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	port := 30053
	caddy.Quiet = true
	dnsserver.Quiet = true
	create(contents)
	defer remove(contents)
	middle := filepath.Join("..", "plugin")
	dirs, err := ioutil.ReadDir(middle)
	if err != nil {
		t.Fatalf("Could not read %s: %q", middle, err)
	}
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		readme := filepath.Join(middle, d.Name())
		readme = filepath.Join(readme, "README.md")
		inputs, err := corefileFromReadme(readme)
		if err != nil {
			continue
		}
		for _, in := range inputs {
			dnsserver.Port = strconv.Itoa(port)
			server, err := caddy.Start(in)
			if err != nil {
				t.Errorf("Failed to start server with %s, for input %q:\n%s", readme, err, in.Body())
			}
			server.Stop()
			port++
		}
	}
}
func corefileFromReadme(readme string) ([]*Input, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	f, err := os.Open(readme)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	input := []*Input{}
	corefile := false
	temp := ""
	for s.Scan() {
		line := s.Text()
		if line == "~~~ corefile" || line == "``` corefile" {
			corefile = true
			continue
		}
		if corefile && (line == "~~~" || line == "```") {
			input = append(input, NewInput(temp))
			temp = ""
			corefile = false
			continue
		}
		if corefile {
			temp += line + "\n"
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return input, nil
}
func create(c map[string]string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for name, content := range c {
		ioutil.WriteFile(name, []byte(content), 0644)
	}
}
func remove(c map[string]string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for name := range c {
		os.Remove(name)
	}
}
