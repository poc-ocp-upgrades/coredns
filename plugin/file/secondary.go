package file

import (
	"math/rand"
	"time"
	"github.com/miekg/dns"
)

func (z *Zone) TransferIn() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(z.TransferFrom) == 0 {
		return nil
	}
	m := new(dns.Msg)
	m.SetAxfr(z.origin)
	z1 := z.CopyWithoutApex()
	var (
		Err	error
		tr	string
	)
Transfer:
	for _, tr = range z.TransferFrom {
		t := new(dns.Transfer)
		c, err := t.In(m, tr)
		if err != nil {
			log.Errorf("Failed to setup transfer `%s' with `%q': %v", z.origin, tr, err)
			Err = err
			continue Transfer
		}
		for env := range c {
			if env.Error != nil {
				log.Errorf("Failed to transfer `%s' from %q: %v", z.origin, tr, env.Error)
				Err = env.Error
				continue Transfer
			}
			for _, rr := range env.RR {
				if err := z1.Insert(rr); err != nil {
					log.Errorf("Failed to parse transfer `%s' from: %q: %v", z.origin, tr, err)
					Err = err
					continue Transfer
				}
			}
		}
		Err = nil
		break
	}
	if Err != nil {
		return Err
	}
	z.Tree = z1.Tree
	z.Apex = z1.Apex
	*z.Expired = false
	log.Infof("Transferred: %s from %s", z.origin, tr)
	return nil
}
func (z *Zone) shouldTransfer() (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := new(dns.Client)
	c.Net = "tcp"
	m := new(dns.Msg)
	m.SetQuestion(z.origin, dns.TypeSOA)
	var Err error
	serial := -1
Transfer:
	for _, tr := range z.TransferFrom {
		Err = nil
		ret, _, err := c.Exchange(m, tr)
		if err != nil || ret.Rcode != dns.RcodeSuccess {
			Err = err
			continue
		}
		for _, a := range ret.Answer {
			if a.Header().Rrtype == dns.TypeSOA {
				serial = int(a.(*dns.SOA).Serial)
				break Transfer
			}
		}
	}
	if serial == -1 {
		return false, Err
	}
	if z.Apex.SOA == nil {
		return true, Err
	}
	return less(z.Apex.SOA.Serial, uint32(serial)), Err
}
func less(a, b uint32) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if a < b {
		return (b - a) <= MaxSerialIncrement
	}
	return (a - b) > MaxSerialIncrement
}
func (z *Zone) Update() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for z.Apex.SOA == nil {
		time.Sleep(1 * time.Second)
	}
	retryActive := false
Restart:
	refresh := time.Second * time.Duration(z.Apex.SOA.Refresh)
	retry := time.Second * time.Duration(z.Apex.SOA.Retry)
	expire := time.Second * time.Duration(z.Apex.SOA.Expire)
	refreshTicker := time.NewTicker(refresh)
	retryTicker := time.NewTicker(retry)
	expireTicker := time.NewTicker(expire)
	for {
		select {
		case <-expireTicker.C:
			if !retryActive {
				break
			}
			*z.Expired = true
		case <-retryTicker.C:
			if !retryActive {
				break
			}
			time.Sleep(jitter(2000))
			ok, err := z.shouldTransfer()
			if err != nil {
				log.Warningf("Failed retry check %s", err)
				continue
			}
			if ok {
				if err := z.TransferIn(); err != nil {
					break
				}
			}
			retryActive = false
			refreshTicker.Stop()
			retryTicker.Stop()
			expireTicker.Stop()
			goto Restart
		case <-refreshTicker.C:
			time.Sleep(jitter(5000))
			ok, err := z.shouldTransfer()
			if err != nil {
				log.Warningf("Failed refresh check %s", err)
				retryActive = true
				continue
			}
			if ok {
				if err := z.TransferIn(); err != nil {
					retryActive = true
					break
				}
			}
			retryActive = false
			refreshTicker.Stop()
			retryTicker.Stop()
			expireTicker.Stop()
			goto Restart
		}
	}
}
func jitter(n int) time.Duration {
	_logClusterCodePath()
	defer _logClusterCodePath()
	r := rand.Intn(n)
	return time.Duration(r) * time.Millisecond
}

const MaxSerialIncrement uint32 = 2147483647
