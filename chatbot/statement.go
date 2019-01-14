package chatbot

import (
	"encoding/json"
	"strings"

	"github.com/epigos/newsbot/utils"

	dgcm "github.com/mlabouardy/dialogflow-go-client/models"
)

// Statement represents a single spoken entity, sentence or
// phrase that someone can say
type Statement struct {
	Text      string        `json:"text"`
	Payload   string        `json:"payload,omitempty"`
	UserID    string        `json:"-"`
	Responses []interface{} `json:"responses"`
	Score     float32       `json:"-"`
	Meta      utils.Map     `json:"meta"`
}

// NewStatement creates and returns a pointer of new Statement
func NewStatement(t, userID string) *Statement {
	s := &Statement{
		Text:   strings.TrimSpace(t),
		UserID: userID,
		Meta:   utils.Map{},
	}
	return s
}

// SerializeResponse statement into JSON
func (s *Statement) SerializeResponse() string {
	bs, _ := json.Marshal(s.Responses)
	return string(bs)
}

// AddResponse for statement
func (s *Statement) AddResponse(r interface{}) {
	s.Responses = append(s.Responses, r)
}

// SetScore set Score of intent
func (s *Statement) SetScore(c float32) {
	s.Score = c
}

// SetPayload set Payload of intent
func (s *Statement) SetPayload(p string) {
	s.Payload = p
}

// AddTextResponse add text response for statement
func (s *Statement) AddTextResponse(text string) {
	s.AddResponse(utils.NewTextMessage(s.UserID, text))
}

func (s *Statement) addMessageResponseFromDialog(ms []dgcm.Message) {
	for _, msg := range ms {
		s.AddResponse(utils.NewTextMessage(s.UserID, msg.Speech))
	}
}
