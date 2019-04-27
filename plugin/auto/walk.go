package auto

import (
	"os"
	"path/filepath"
	"regexp"
	"github.com/coredns/coredns/plugin/file"
	"github.com/miekg/dns"
)

func (a Auto) Walk() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	toDelete := make(map[string]bool)
	for _, n := range a.Zones.Names() {
		toDelete[n] = true
	}
	filepath.Walk(a.loader.directory, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		match, origin := matches(a.loader.re, info.Name(), a.loader.template)
		if !match {
			return nil
		}
		if z, ok := a.Zones.Z[origin]; ok {
			toDelete[origin] = false
			z.SetFile(path)
			return nil
		}
		reader, err := os.Open(path)
		if err != nil {
			log.Warningf("Opening %s failed: %s", path, err)
			return nil
		}
		defer reader.Close()
		zo, err := file.Parse(reader, origin, path, 0)
		if err != nil {
			log.Warningf("Parse zone `%s': %v", origin, err)
			return nil
		}
		zo.ReloadInterval = a.loader.ReloadInterval
		zo.Upstream = a.loader.upstream
		zo.TransferTo = a.loader.transferTo
		a.Zones.Add(zo, origin)
		if a.metrics != nil {
			a.metrics.AddZone(origin)
		}
		zo.Notify()
		log.Infof("Inserting zone `%s' from: %s", origin, path)
		toDelete[origin] = false
		return nil
	})
	for origin, ok := range toDelete {
		if !ok {
			continue
		}
		if a.metrics != nil {
			a.metrics.RemoveZone(origin)
		}
		a.Zones.Remove(origin)
		log.Infof("Deleting zone `%s'", origin)
	}
	return nil
}
func matches(re *regexp.Regexp, filename, template string) (match bool, origin string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	base := filepath.Base(filename)
	matches := re.FindStringSubmatchIndex(base)
	if matches == nil {
		return false, ""
	}
	by := re.ExpandString(nil, template, base, matches)
	if by == nil {
		return false, ""
	}
	origin = dns.Fqdn(string(by))
	return true, origin
}
