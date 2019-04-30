package file

import (
	"os"
	"time"
)

var TickTime = 1 * time.Second

func (z *Zone) Reload() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if z.ReloadInterval == 0 {
		return nil
	}
	tick := time.NewTicker(TickTime)
	go func() {
		for {
			select {
			case <-tick.C:
				if z.LastReloaded.Add(z.ReloadInterval).After(time.Now()) {
					continue
				}
				z.LastReloaded = time.Now()
				zFile := z.File()
				reader, err := os.Open(zFile)
				if err != nil {
					log.Errorf("Failed to open zone %q in %q: %v", z.origin, zFile, err)
					continue
				}
				serial := z.SOASerialIfDefined()
				zone, err := Parse(reader, z.origin, zFile, serial)
				if err != nil {
					if _, ok := err.(*serialErr); !ok {
						log.Errorf("Parsing zone %q: %v", z.origin, err)
					}
					continue
				}
				z.reloadMu.Lock()
				z.Apex = zone.Apex
				z.Tree = zone.Tree
				z.reloadMu.Unlock()
				log.Infof("Successfully reloaded zone %q in %q with serial %d", z.origin, zFile, z.Apex.SOA.Serial)
				z.Notify()
			case <-z.reloadShutdown:
				tick.Stop()
				return
			}
		}
	}()
	return nil
}
func (z *Zone) SOASerialIfDefined() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	z.reloadMu.Lock()
	defer z.reloadMu.Unlock()
	if z.Apex.SOA != nil {
		return int64(z.Apex.SOA.Serial)
	}
	return -1
}
