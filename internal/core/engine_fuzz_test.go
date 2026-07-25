package core

import (
	"testing"
)

func FuzzEngineScan(f *testing.F) {
	f.Add("192.168.1.1")
	f.Add("example.com")
	f.Add("web01")
	f.Add("::1")
	f.Fuzz(func(t *testing.T, target string) {
		e := NewEngine()
		_, _ = e.Scan(target)
	})
}
