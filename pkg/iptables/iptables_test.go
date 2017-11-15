package iptables

import (
	"testing"
)

func TestListChains(t *testing.T) {
	ipt := New()
	_, err := ipt.ListChains()
	if err != nil {
		t.Errorf("Res: %v", err)
		t.Fail()
	}
}

func TestStats(t *testing.T) {
	ipt := New()
	err := ipt.Stats()
	if err != nil {
		t.Errorf("Res: %v", err)
		t.Fail()
	}
}
