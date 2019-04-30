package test

import (
	"io/ioutil"
	"os"
)

func TempFile(dir, content string) (string, func(), error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f, err := ioutil.TempFile(dir, "go-test-tmpfile")
	if err != nil {
		return "", nil, err
	}
	if err := ioutil.WriteFile(f.Name(), []byte(content), 0644); err != nil {
		return "", nil, err
	}
	rmFunc := func() {
		os.Remove(f.Name())
	}
	return f.Name(), rmFunc, nil
}
