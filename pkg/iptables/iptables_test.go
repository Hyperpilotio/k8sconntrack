package iptables

import (
    "testing"
    "fmt"
)

func TestGather(t *testing.T) {
    ipt := NewIptables()
    ipt.Chains = []string{"INPUT", "OUTPUT"}
    ipt.Table = "filter"
    err := ipt.Gather()
    if err != nil {
        fmt.Printf("Res: %v", err)
        t.Fail()
    }
}

func TestParseAndGather(t *testing.T) {
    data := []string{
        `Chain INPUT (policy ACCEPT 11179 packets, 607687 bytes)
    pkts      bytes target     prot opt in     out     source               destination
     203    15347 ACCEPT     all  --  lo     *       0.0.0.0/0            0.0.0.0/0
       0        0 ACCEPT     all  --  lo     *       0.0.0.0/0            0.0.0.0/0
       0        0 ACCEPT     all  --  lo     *       0.0.0.0/0            0.0.0.0/0
       0        0 ACCEPT     all  --  lo     *       0.0.0.0/0            0.0.0.0/0
       0        0 ACCEPT     all  --  lo     *       0.0.0.0/0            0.0.0.0/0
      65     2848 DROP       all  --  *      *       104.31.88.86         0.0.0.0/0`,
        `Chain INPUT (policy DROP 0 packets, 0 bytes)
    pkts bytes target     prot opt in     out     source               destination
    100   1024   ACCEPT     tcp  --  *      *       192.168.0.0/24       0.0.0.0/0            tcp dpt:22 /* ssh */
     42   2048   ACCEPT     tcp  --  *      *       192.168.0.0/24       0.0.0.0/0            tcp dpt:80 /* httpd */`,
    }

    ipt := NewIptables()
    for _, val := range data {
        err := ipt.parseAndGather(val)
        if err != nil {
            fmt.Printf("Res: %v", err)
            t.Fail()
        }
    }

}

func TestChainList(t *testing.T) {
    ipt := NewIptables()
    _, err := ipt.chainList("filter", "INPUT")
    if err != nil {
        fmt.Printf("Res: %v", err)
        t.Fail()
    }
}

