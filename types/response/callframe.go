package response

// https://github.com/ethereum/go-ethereum/blob/ad15050c7fbedd0f05a49e81400de18c2cc2c284/eth/tracers/native/call.go
type CallFrame struct {
	Type    string      `json:"type"`
	From    string      `json:"from"`
	To      string      `json:"to,omitempty"`
	Value   string      `json:"value,omitempty"`
	Gas     string      `json:"gas"`
	GasUsed string      `json:"gasUsed"`
	Input   string      `json:"input"`
	Output  string      `json:"output,omitempty"`
	Error   string      `json:"error,omitempty"`
	Calls   []CallFrame `json:"calls,omitempty"`
}
