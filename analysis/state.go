package analysis

import (
	"cwrenhold/lsp-from-scratch/lsp"
	"fmt"
)

type State struct {
	// Map of file names to their content
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: make(map[string]string),
	}
}

func (s *State) OpenDocument(uri string, content string) {
	s.Documents[uri] = content
}

func (s *State) UpdateDocument(uri string, content string) {
	s.Documents[uri] = content
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	// In the real world, we'd look up the type or something like that

	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
		},
	}
}
