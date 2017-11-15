package transactioncounter

import (
	"testing"
)

func TestCount(t *testing.T) {
	tests := []struct {
		CounterMap      map[string]map[string]EndpointAbsCounter
		ServiceName     string
		EndpointAddress string
		ExpectedCount   uint64
		ExpectedBytes   uint64
		ExpectedPackets uint64
	}{
		{
			CounterMap:      map[string]map[string]EndpointAbsCounter{},
			ServiceName:     "service1",
			EndpointAddress: "10.0.0.2",
			ExpectedCount:   1,
			ExpectedBytes:   32,
			ExpectedPackets: 1,
		},
		{
			CounterMap: map[string]map[string]EndpointAbsCounter{
				"service1": map[string]EndpointAbsCounter{
					"10.0.1.2": EndpointAbsCounter{Counts: 3},
				},
			},
			ServiceName:     "service1",
			EndpointAddress: "10.0.0.2",
			ExpectedCount:   1,
			ExpectedBytes:   32,
			ExpectedPackets: 1,
		},
		{
			CounterMap: map[string]map[string]EndpointAbsCounter{
				"service1": map[string]EndpointAbsCounter{
					"10.0.0.2": EndpointAbsCounter{Counts: 3},
				},
			},
			ServiceName:     "service1",
			EndpointAddress: "10.0.0.2",
			ExpectedCount:   4,
			ExpectedBytes:   32,
			ExpectedPackets: 1,
		},
	}

	for _, test := range tests {
		transactionCounter := NewTransactionCounter(nil)
		transactionCounter.counter = test.CounterMap
		transactionCounter.Count([]*countInfo{&countInfo{serviceName: test.ServiceName, endpointAddress: test.EndpointAddress}})
		var c uint64
		if epMap, exist := transactionCounter.counter[test.ServiceName]; exist {
			if epCounts, has := epMap[test.EndpointAddress]; has {
				c = epCounts.Counts
			} else {
				t.Errorf("Endpoint %s is not found in transaction info for service %s.", test.EndpointAddress, test.ServiceName)
			}
		} else {
			t.Errorf("Service %s is not found in transaction counter map.", test.ServiceName)
		}
		if c != test.ExpectedCount {
			t.Errorf("Expected count is %d, got %d", test.ExpectedCount, c)
		}
	}
}
