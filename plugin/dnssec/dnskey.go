package dnssec

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"os"
	"time"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type DNSKEY struct {
	K	*dns.DNSKEY
	D	*dns.DS
	s	crypto.Signer
	tag	uint16
}

func ParseKeyFile(pubFile, privFile string) (*DNSKEY, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	f, e := os.Open(pubFile)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	k, e := dns.ReadRR(f, pubFile)
	if e != nil {
		return nil, e
	}
	f, e = os.Open(privFile)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	dk, ok := k.(*dns.DNSKEY)
	if !ok {
		return nil, errors.New("no public key found")
	}
	p, e := dk.ReadPrivateKey(f, privFile)
	if e != nil {
		return nil, e
	}
	if s, ok := p.(*rsa.PrivateKey); ok {
		return &DNSKEY{K: dk, D: dk.ToDS(dns.SHA256), s: s, tag: dk.KeyTag()}, nil
	}
	if s, ok := p.(*ecdsa.PrivateKey); ok {
		return &DNSKEY{K: dk, D: dk.ToDS(dns.SHA256), s: s, tag: dk.KeyTag()}, nil
	}
	return &DNSKEY{K: dk, D: dk.ToDS(dns.SHA256), s: nil, tag: 0}, errors.New("no private key found")
}
func (d Dnssec) getDNSKEY(state request.Request, zone string, do bool, server string) *dns.Msg {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	keys := make([]dns.RR, len(d.keys))
	for i, k := range d.keys {
		keys[i] = dns.Copy(k.K)
		keys[i].Header().Name = zone
	}
	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Answer = keys
	if !do {
		return m
	}
	incep, expir := incepExpir(time.Now().UTC())
	if sigs, err := d.sign(keys, zone, 3600, incep, expir, server); err == nil {
		m.Answer = append(m.Answer, sigs...)
	}
	return m
}
func (k DNSKEY) isZSK() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return k.K.Flags&(1<<8) == (1<<8) && k.K.Flags&1 == 0
}
func (k DNSKEY) isKSK() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return k.K.Flags&(1<<8) == (1<<8) && k.K.Flags&1 == 1
}
