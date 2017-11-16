package transactioncounter

import (
	"fmt"
	"sync"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/types"

	"github.com/Hyperpilotio/k8sconntrack/pkg/conntrack"

	"github.com/golang/glog"
)

type endpointsInfo struct {
	types.NamespacedName
}

type TransactionCounter struct {
	conntrack *conntrack.ConnTrack

	mu sync.Mutex

	endpointsMap map[string]*endpointsInfo

	// key is service name, value is the transaction related to it.
	counter map[string]map[string]EndpointAbsCounter

	lastPollTimestamp uint64
}

func NewTransactionCounter(conntrack *conntrack.ConnTrack) *TransactionCounter {
	return &TransactionCounter{
		counter:   make(map[string]map[string]EndpointAbsCounter),
		conntrack: conntrack,

		endpointsMap: make(map[string]*endpointsInfo),
	}
}

// Implement k8s.io/pkg/proxy/config/EndpointsConfigHandler Interface.
func (this *TransactionCounter) OnEndpointsUpdate(allEndpoints []api.Endpoints) {
	start := time.Now()
	defer func() {
		glog.V(4).Infof("OnEndpointsUpdate took %v for %d endpoints", time.Since(start), len(allEndpoints))
	}()

	this.mu.Lock()
	defer this.mu.Unlock()

	// Clear the current endpoints set.
	this.endpointsMap = make(map[string]*endpointsInfo)

	for i := range allEndpoints {
		endpoints := &allEndpoints[i]
		for j := range endpoints.Subsets {
			ss := &endpoints.Subsets[j]
			for k := range ss.Addresses {
				addr := &ss.Addresses[k]
				this.endpointsMap[addr.IP] = &endpointsInfo{types.NamespacedName{Namespace: endpoints.Namespace, Name: endpoints.Name}}
			}
		}
	}

	this.syncConntrack()
}

// Clear the transaction counter map.
func (tc *TransactionCounter) Reset() {
	glog.V(3).Infof("Inside reset transaction counter")
	counterMap := make(map[string]map[string]EndpointAbsCounter)

	tc.counter = counterMap

	// As after each poll, the counter map is cleaned, so this is the right place to set the lastPollTimestamp.
	tc.lastPollTimestamp = uint64(time.Now().Unix())
}

// Increment the transaction count for a single endpoint.
// Transaction counter map uses serviceName as key and endpoint map as value.
// In endpoint map, key is endpoint IP address, value is the number of transaction happened on the endpoint.
func (tc *TransactionCounter) Count(infos []*countInfo) {
	for _, info := range infos {
		serviceName := info.serviceName
		endpointAddress := info.endpointAddress
		epMap, ok := tc.counter[serviceName]
		if !ok {
			glog.V(4).Infof("Service %s is not tracked. Now initializing in map", serviceName)

			epMap = make(map[string]EndpointAbsCounter)
		}
		// count, ok := epMap[endpointAddress]
		epCollector, ok := epMap[endpointAddress]
		if !ok {
			glog.V(4).Infof("Endpoint %s for Service %s is not tracked. Now initializing in map", endpointAddress, serviceName)
			epCollector = EndpointAbsCounter{}
		}
		epCollector = EndpointAbsCounter{
			Counts:  epCollector.Counts + 1,
			Bytes:   epCollector.Bytes + info.bytes,
			Packets: epCollector.Packets + info.packets,
			Role:    info.role,
		}
		epMap[endpointAddress] = epCollector
		tc.counter[serviceName] = epMap
		glog.V(4).Infof("Transaction count of %s is %+v", endpointAddress, epMap[endpointAddress])
	}
}

func (tc *TransactionCounter) GetAllTransactions() []*Transaction {
	var transactions []*Transaction

	// Here we need to translate the absolute count value into count/second.
	if tc.lastPollTimestamp == 0 {
		// When lastPollTimestamp is 0, meaning that current poll is the first poll. We cannot get the count/s, so just return.
		return transactions
	}
	// Get the time difference between two poll.
	timeDiff := uint64(time.Now().Unix()) - tc.lastPollTimestamp
	glog.V(4).Infof("Time diff is %d", timeDiff)

	for svcName, epMap := range tc.counter {
		// Before append, change count to count per second.
		valueMap := make(map[string]EndpointCounter)
		countMap := make(map[string]EndpointAbsCounter)
		for ep, epCounter := range epMap {
			value := EndpointCounter{
				Counts:  float64(epCounter.Counts) / float64(timeDiff),
				Bytes:   float64(epCounter.Bytes) / float64(timeDiff),
				Packets: float64(epCounter.Packets) / float64(timeDiff),
				Role:    epCounter.Role,
			}
			valueMap[ep] = value

			count := EndpointAbsCounter{
				Counts:  epCounter.Counts,
				Bytes:   epCounter.Bytes,
				Packets: epCounter.Packets,
				Role:    epCounter.Role,
			}
			countMap[ep] = count
		}
		transaction := &Transaction{
			ServiceId:           svcName,
			EndpointsCounterMap: valueMap,
			EpCountAbs:          countMap,
		}
		glog.V(1).Infof("Get transaction data: %++v", transaction)
		transactions = append(transactions, transaction)
	}

	return transactions
}

// Get all the current Established TCP connections from conntrack and add count to transaction counter.
func (this *TransactionCounter) ProcessConntrackConnections() {
	this.mu.Lock()
	defer this.mu.Unlock()

	this.syncConntrack()
}

func (this *TransactionCounter) syncConntrack() {
	connections := this.conntrack.ConnectionEvents()
	if len(connections) > 0 {
		glog.V(3).Infof("Connections:\n")
		for _, cn := range connections {
			infos := this.preProcessConnections(cn)
			this.Count(infos)
		}
	}
}

type countInfo struct {
	serviceName     string
	endpointAddress string
	bytes           uint64
	packets         uint64
	role            string
}

// Filter out connection does not have endpoints address as either Local or Remote Address
func (this *TransactionCounter) preProcessConnections(c conntrack.ConntrackInfo) []*countInfo {
	var infos []*countInfo
	if svcName, exist := this.endpointsMap[c.Src.String()]; exist {
		infos = append(infos, &countInfo{svcName.String(), fmt.Sprintf("%s:%v", c.Src.String(), c.SrcPort), c.Bytes, c.Packets, "sender"})
	}
	if svcName, exist := this.endpointsMap[c.Dst.String()]; exist {
		infos = append(infos, &countInfo{svcName.String(), fmt.Sprintf("%s:%v", c.Dst.String(), c.DstPort), c.Bytes, c.Packets, "receiver"})
	}
	return infos

}
