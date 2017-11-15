package transactioncounter

type EndpointCounter struct {
	Counts  float64 `json:"count,omitempty"`
	Bytes   float64 `json:"bytes,omitempty"`
	Packets float64 `json:"packets,omitempty"`
	Role    string  `json:"role,omitempty"`
}

type EndpointAbsCounter struct {
	Counts  uint64 `json:"count,omitempty"`
	Bytes   uint64 `json:"bytes,omitempty"`
	Packets uint64 `json:"packets,omitempty"`
	Role    string `json:"role,omitempty"`
}

type Transaction struct {
	ServiceId           string                        `json:"serviceID,omitempty"`
	EndpointsCounterMap map[string]EndpointCounter    `json:"endpointCounter,omitempty"`
	EpCountAbs          map[string]EndpointAbsCounter `json:"endpointAbs,omitempty"`
}

func (this *Transaction) GetEndpointsCounterMap() map[string]EndpointCounter {
	return this.EndpointsCounterMap
}
