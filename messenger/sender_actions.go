package messenger

const (
	// ActionMarkSeen action for mark as seen
	ActionMarkSeen = "mark_seen"
	// ActionTypingOn action for typing on
	ActionTypingOn = "typing_on"
	// ActionTypingOff action for typing off
	ActionTypingOff = "typing_off"
)

// SenderAction sends facebook sender action
type SenderAction struct {
	Recipient *Recipient `json:"recipient"`
	Action    string     `json:"sender_action"`
}

// NewSenderAction returns a pointer of new SenderAction
func NewSenderAction(userID string, s string) *SenderAction {
	return &SenderAction{&Recipient{ID: userID}, s}
}

func (mg *Messenger) sendAction(s *Recipient, action string) {
	a := NewSenderAction(s.ID, action)
	resp, err := mg.makeFbRequest(messagesPath, "POST", a)
	if err != nil {
		logger.Error(err)
	}
	_, err = mg.decodeResponse(resp)
	if err != nil {
		logger.Error(err)
	}
}

// MarkSeen mark facebook message as seen
func (mg *Messenger) MarkSeen(s *Recipient) {
	mg.sendAction(s, ActionMarkSeen)
}

// SendTypingOn show typing on to user
func (mg *Messenger) SendTypingOn(s *Recipient) {
	mg.sendAction(s, ActionTypingOn)
}

// SendTypingOff show typing on to user
func (mg *Messenger) SendTypingOff(s *Recipient) {
	mg.sendAction(s, ActionTypingOff)
}
