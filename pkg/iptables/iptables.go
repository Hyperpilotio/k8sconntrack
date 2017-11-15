package iptables

import (
    "fmt"

    "github.com/golang/glog"

    ipt "github.com/coreos/go-iptables/iptables"
)

var (
    t *ipt.IPTables
    c *Collector
    defaultTables = []string{"filter", "nat", "mangle", "raw"}
)

func init() {
    var err error
    t, err = ipt.New()
    if err != nil {
        // TODO init glog at other place
        glog.Warningf("failed to init go-iptables err:%+v", err)
    }
    c = New()
}

type Chain struct {
    Data [][]string     `json:"data"`
    Name string         `json:"name"`
}

type Table struct {
    Chains  []*Chain     `json:"chains"`
    Name    string      `json:"name"`
}

type Collector struct {
    Tables  []Table
}

func New() *Collector {
    tables := []Table{}
    for _, table := range defaultTables {
        chainList, err := t.ListChains(table)
        if err != nil {

        }
        chains := []*Chain{}
        for _, chainName := range chainList {
            chains = append(chains, &Chain{Name: chainName})
        }
        tables = append(tables, Table{Name: table, Chains: chains})
    }

    c := &Collector{tables}
    return c
}

func (c *Collector) ListChains() (map[string][]string, error) {
    tableMap := map[string][]string{}
    for _, table := range c.Tables {
        chains, err := t.ListChains(table.Name)
        if err != nil {
            glog.Warningf(
                "Unable to list chains while table=%s err=%s",
                table.Name, err.Error())
            return nil, fmt.Errorf("Unable to list chains while table=%s err=%s", table.Name, err.Error())
        }
        tableMap[table.Name] = chains
    }
    return tableMap, nil
}

func (c *Collector) String() string {
    return fmt.Sprintf("%+v", *c)
}

func (c *Collector) Stats() error {
    for _, table := range c.Tables {
        for _, chain := range table.Chains {
            data, err := t.Stats(table.Name, chain.Name)
            if err != nil {
                glog.Warningf(
                    "Unable to collect data from iptables while table=%s chain=%s err=%s",
                    table.Name, chain.Name, err.Error())
                continue
            }
            chain.Data = data
        }
    }
    return nil
}

