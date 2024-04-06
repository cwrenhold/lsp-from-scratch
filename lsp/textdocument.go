package lsp

type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
	// Ignore the parameters that we're not supporting yet
	// Range       Range  `json:"range"`
	// RangeLength int    `json:"rangeLength"`
}
