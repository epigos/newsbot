package chatbot

// BestLogic best logic adapter
type BestLogic struct {
	bot    *Chatbot
	logics []LogicAdapter
}

// NewBestLogic returns new BestLogic
func NewBestLogic() *BestLogic {
	return &BestLogic{
		logics: []LogicAdapter{
			NewPostBackLogic(),
			NewDialogFlowLogic(),
		},
	}
}

func (l *BestLogic) setChatbot(b *Chatbot) {
	l.bot = b
	for _, logic := range l.logics {
		logic.setChatbot(b)
	}
}

func (l *BestLogic) canProcess(s *Statement) bool {
	return true
}

// Process reads the user's input from the terminal.
func (l *BestLogic) Process(st *Statement) *Statement {
	l.bot.Logger.Info("Finding best response...")

	var response *Statement

	for _, logic := range l.logics {
		if logic.canProcess(st) {
			response = logic.Process(st)
			break
		}
	}
	return response
}
