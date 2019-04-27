package file

import (
	"fmt"
	"strings"
)

func ExampleZone_All() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	zone, err := Parse(strings.NewReader(dbMiekNL), testzone, "stdin", 0)
	if err != nil {
		return
	}
	records := zone.All()
	for _, r := range records {
		fmt.Printf("%+v\n", r)
	}
}
