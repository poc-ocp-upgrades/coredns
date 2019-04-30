package log

import (
	"fmt"
	"os"
)

type P struct{ plugin string }

func NewWithPlugin(name string) P {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return P{"plugin/" + name + ": "}
}
func (p P) logf(level, format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	log(level, p.plugin, fmt.Sprintf(format, v...))
}
func (p P) log(level string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	log(level+p.plugin, v...)
}
func (p P) Debug(v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !D {
		return
	}
	p.log(debug, v...)
}
func (p P) Debugf(format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !D {
		return
	}
	p.logf(debug, format, v...)
}
func (p P) Info(v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.log(info, v...)
}
func (p P) Infof(format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.logf(info, format, v...)
}
func (p P) Warning(v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.log(warning, v...)
}
func (p P) Warningf(format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.logf(warning, format, v...)
}
func (p P) Error(v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.log(err, v...)
}
func (p P) Errorf(format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.logf(err, format, v...)
}
func (p P) Fatal(v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.log(fatal, v...)
	os.Exit(1)
}
func (p P) Fatalf(format string, v ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	p.logf(fatal, format, v...)
	os.Exit(1)
}
