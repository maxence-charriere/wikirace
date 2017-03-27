package main

import "testing"

func TestEpureTarget(t *testing.T) {
	if res := epureTarget("   Mike Tyson  "); res != "Mike_Tyson" {
		t.Errorf(`res should be "Mike_Tyson": "%v"`, res)
	}
}
