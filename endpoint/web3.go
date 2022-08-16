package endpoint

type Web3 struct {
	*Endpoint
}

func NewWeb3(endpoint *Endpoint) *Web3 {
	return &Web3{endpoint}
}

func (ep *Web3) ClientVersion() string {
	return "Aurora"
}
