package reload

import (
	"crypto/md5"
	"time"
	"github.com/mholt/caddy"
)

const (
	unused		= 0
	maybeUsed	= 1
	used		= 2
)

type reload struct {
	interval	time.Duration
	usage		int
	quit		chan bool
}

func hook(event caddy.EventName, info interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if event != caddy.InstanceStartupEvent {
		return nil
	}
	if r.usage == unused {
		return nil
	}
	instance := info.(*caddy.Instance)
	md5sum := md5.Sum(instance.Caddyfile().Body())
	log.Infof("Running configuration MD5 = %x\n", md5sum)
	go func() {
		tick := time.NewTicker(r.interval)
		for {
			select {
			case <-tick.C:
				corefile, err := caddy.LoadCaddyfile(instance.Caddyfile().ServerType())
				if err != nil {
					continue
				}
				s := md5.Sum(corefile.Body())
				if s != md5sum {
					md5sum = s
					r.usage = maybeUsed
					_, err := instance.Restart(corefile)
					if err != nil {
						log.Errorf("Corefile changed but reload failed: %s\n", err)
						continue
					}
					if r.usage == maybeUsed {
						r.usage = unused
					}
					return
				}
			case <-r.quit:
				return
			}
		}
	}()
	return nil
}
