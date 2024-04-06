package analysis

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
