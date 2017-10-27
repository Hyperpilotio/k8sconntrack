package iptables

import (
    "fmt"

    "github.com/golang/glog"

    ipt "github.com/coreos/go-iptables/iptables"
)

var t *ipt.IPTables
var c *Collector

func init() {
    var err error
    t, err = ipt.New()
    if err != nil {
        // TODO init glog at other place
        glog.Warningf("failed to init go-iptables err:%+v", err)
    }
    c = New("", []string{""})
}

type Chain struct {
    Data [][]string
    Name string
}

type Collector struct {
    Chains  []Chain

    Table  string
}

func New(table string, chains []string) *Collector {
    c := &Collector{Table: table, Chains: []Chain{}}
    for _, val := range chains {
        c.Chains = append(c.Chains, Chain{Name: val})
    }
    return c
}

func (c *Collector) Stats() error {
    for index := range c.Chains {
        data, err := t.Stats(c.Table, c.Chains[index].Name)
        if err != nil {
            glog.Warningf(
                "Unable to collect data from iptables while table=%s chain=%s err=%s",
                c.Table, c.Chains[index].Name, err.Error())
            continue
        }
        c.Chains[index].Data = data
    }
    return nil
}

func (c *Collector) ToString() string {
    return fmt.Sprintf("%+v", *c)
}

