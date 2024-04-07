package analysis

import (
	"cwrenhold/lsp-from-scratch/lsp"
	"fmt"
	"strings"
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

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}

	for row, line := range strings.Split(text, "\n") {
		searchText := "VS Code"
		if strings.Contains(line, searchText) {
			index := strings.Index(line, searchText)
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range: LineRange(row, index, index+len(searchText)),
				Severity: 1,
				Source:   "Common Sense",
				Message:  "VS Code is a bad time, so sad",
			})
		}

		searchText = "neovim"
		if strings.Contains(line, searchText) {
			index := strings.Index(line, searchText)
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range: LineRange(row, index, index+len(searchText)),
				Severity: 4,
				Source:   "Common Sense",
				Message:  "Neovim is a good time :)",
			})
		}
	}

	return diagnostics
}

func (s *State) OpenDocument(uri string, content string) []lsp.Diagnostic {
	s.Documents[uri] = content

	return getDiagnosticsForFile(content)
}

func (s *State) UpdateDocument(uri string, content string) []lsp.Diagnostic {
	s.Documents[uri] = content

	return getDiagnosticsForFile(content)
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

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	// In the real world, we'd look up the location of the definition

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line: position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line: position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) TextDocumentCodeAction(id int, uri string) lsp.TextDocumentCodeActionResponse {
	text := s.Documents[uri]

	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		searchText := "VS Code"
		idx := strings.Index(line, searchText)
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range: LineRange(row, idx, idx+len(searchText)),
					NewText: "Neovim",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: fmt.Sprintf("Replace %s with Neovim", searchText),
				Edit: &lsp.WorkspaceEdit{
					Changes: replaceChange,
				},
			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[uri] = []lsp.TextEdit{
				{
					Range: LineRange(row, idx, idx+len(searchText)),
					NewText: strings.Repeat("*", len(searchText)),
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: fmt.Sprintf("Censor %s", searchText),
				Edit: &lsp.WorkspaceEdit{
					Changes: censorChange,
				},
			})
		}
	}

	response := lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}

	return response
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
	items := []lsp.CompletionItem{
		{
			Label: "Neovim",
			Detail: "Text editor",
			Documentation: "Neovim is a text editor",
		},
	}

	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}

	return response
}

func LineRange(line int, start int, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
