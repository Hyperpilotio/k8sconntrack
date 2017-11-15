package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	fcollector "github.com/Hyperpilotio/k8sconntrack/pkg/flowcollector"
	iptables "github.com/Hyperpilotio/k8sconntrack/pkg/iptables"
	tcounter "github.com/Hyperpilotio/k8sconntrack/pkg/transactioncounter"

	"github.com/golang/glog"
)

// Server is a http.Handler which exposes kubelet functionality over HTTP.
type Server struct {
	counter           *tcounter.TransactionCounter
	flowCollector     *fcollector.FlowCollector
	iptablesCollector *iptables.Collector
	mux               *http.ServeMux
}

// NewServer initializes and configures a kubelet.Server object to handle HTTP requests.
func NewServer(counter *tcounter.TransactionCounter, flowCollector *fcollector.FlowCollector, iptablesCollector *iptables.Collector) Server {
	server := Server{
		counter:           counter,
		flowCollector:     flowCollector,
		iptablesCollector: iptablesCollector,
		mux:               http.NewServeMux(),
	}
	server.InstallDefaultHandlers()
	return server
}

// InstallDefaultHandlers registers the default set of supported HTTP request patterns with the mux.
func (s *Server) InstallDefaultHandlers() {
	s.mux.HandleFunc("/", handler)
	s.mux.HandleFunc("/transactions/count", s.getTransactionsCount)
	s.mux.HandleFunc("/transactions", s.getAllTransactionsAndReset)
	s.mux.HandleFunc("/chains", s.getAllChains)
	s.mux.HandleFunc("/flows", s.getAllFlows)
	s.mux.HandleFunc("/iptables", s.getIptables)
	s.mux.HandleFunc("/iptables/chains", s.getAllChains)
}

// ServeHTTP responds to HTTP requests on the Kubelet.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.mux.ServeHTTP(w, req)
}

func (s *Server) getAllTransactionsAndReset(w http.ResponseWriter, r *http.Request) {
	transactions := s.counter.GetAllTransactions()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		panic(err)
	}

	s.resetCounter()
}

func (s *Server) getTransactionsCount(w http.ResponseWriter, r *http.Request) {
	if s.counter == nil {
		fmt.Fprintf(w, "Connection Counter is disabled.")
		return
	}
	transactions := s.counter.GetAllTransactions()

	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (s *Server) getAllFlows(w http.ResponseWriter, r *http.Request) {
	if s.flowCollector == nil {
		fmt.Fprintf(w, "Flow Collector is disabled.")
		return
	}
	flows := s.flowCollector.GetAllFlows()

	data, err := json.MarshalIndent(flows, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (s *Server) resetCounter() {
	s.counter.Reset()
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Vmturbo k8sconntrack Service.")
}

// TODO: For now the address and port number is hardcoded. The actual port number need to be discussed.
func ListenAndServeProxyServer(bindAddress, bindPort string, counter *tcounter.TransactionCounter, flowCollector *fcollector.FlowCollector, iptablesCollector *iptables.Collector) {
	glog.V(3).Infof("Start VMT Kube-proxy server")
	handler := NewServer(counter, flowCollector, iptablesCollector)
	s := &http.Server{
		Addr:           net.JoinHostPort(bindAddress, bindPort),
		Handler:        &handler,
		MaxHeaderBytes: 1 << 20,
	}
	glog.Fatal(s.ListenAndServe())
}

func (s *Server) getIptables(w http.ResponseWriter, r *http.Request) {
	if s.iptablesCollector == nil {
		fmt.Fprintf(w, "Iptables Collector is disabled.")
		return
	}
	err := s.iptablesCollector.Stats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

    tableMap := map[string]iptables.Table{}
    for _, table := range s.iptablesCollector.Tables {
        tableMap[table.Name] = table
    }

	data, err := json.MarshalIndent(tableMap, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (s *Server) getAllChains(w http.ResponseWriter, r *http.Request) {
	if s.iptablesCollector == nil {
		fmt.Fprintf(w, "Iptables Collector is disabled.")
		return
	}
	chains, err := s.iptablesCollector.ListChains()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data, err := json.MarshalIndent(chains, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
