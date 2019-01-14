package chatbot

import (
	"github.com/epigos/newsbot/utils"
	"github.com/epigos/newsbot/web"
)

// Chatbot A convensational chat dialog
type Chatbot struct {
	Name   string
	Logger *utils.Logger
	Logic  LogicAdapter
}

// New creates a new pointer of Chatbot
func New(n string) *Chatbot {
	bot := &Chatbot{
		Name:   n,
		Logger: utils.NewLogger(n),
		Logic:  NewBestLogic(),
	}
	// initialize bot functions
	bot.initialize()
	return bot
}

func (b *Chatbot) String() string {
	return b.Name
}

// Initialize initialize chatbot functions
func (b *Chatbot) initialize() {
	b.Logic.setChatbot(b)
	b.Logger.Info("Initialized bot:", b)
}

// GetResponse generates a response for the input text
func (b *Chatbot) GetResponse(st *Statement) *Statement {
	// get response statement
	response := b.Logic.Process(st)
	// return output
	return response
}

// TestHandler test handler for chatbot
func (b *Chatbot) TestHandler(ctx *web.Context) *web.HTTPError {
	text := ctx.PostValues().Get("text", "")
	if text == "" {
		return ctx.BadRequest("Text is required")
	}

	st := NewStatement(text.(string), "1403078893046594")
	st.SetPayload(ctx.PostValues().Get("payload", "").(string))

	response := b.GetResponse(st)
	return ctx.WriteJSON(response)
}
