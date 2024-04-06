package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`

	// We'll ignore params for now, and specify this in the Request types later
}

type Response struct {
	RPC string `json:"jsonrpc"`
	ID  *int   `json:"id,omitempty"`

	// Result and Error will be handled with subtypes
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}
