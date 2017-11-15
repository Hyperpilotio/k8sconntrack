package iptables

import (
    "testing"
    "fmt"
)

func TestListChains(t *testing.T) {
    ipt := New()
    chains, err := ipt.ListChains()
    if err != nil {
        fmt.Printf("Res: %v", err)
        t.Fail()
    }
    fmt.Println(chains)
}

func TestStats(t *testing.T) {
    ipt := New()
    err := ipt.Stats()
    if err != nil {
        fmt.Printf("Res: %v", err)
        t.Fail()
    }
    fmt.Printf("%+v", ipt.Tables)
}

