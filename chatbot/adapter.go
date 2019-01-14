package chatbot

// Adapter an interface for all adapters.
type Adapter interface {
	setChatbot(b *Chatbot) // Gives the adapter access to the chatbot pointer.
}

// LogicAdapter an interface that represents
// all logic adapters should implement.
type LogicAdapter interface {
	Adapter
	canProcess(s *Statement) bool
	Process(s *Statement) *Statement
}
