package iptables

import (
    "testing"
    "fmt"
)

func TestGather(t *testing.T) {
    ipt := New("filter", []string{"INPUT"})
    err := ipt.Stats()
    if err != nil {
        fmt.Printf("Res: %v", err)
        t.Fail()
    }
    fmt.Println(ipt)
}

