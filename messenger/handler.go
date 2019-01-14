package messenger

import (
	"github.com/epigos/newsbot/chatbot"
	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/utils"
)

// MessageHandler an interface for messenger message handlers
type MessageHandler interface {
	ProcessMessage(m *messaging)
	ProcessDelivery(m *messaging)
	ProcessPostback(m *messaging)
}

// DefaultHandler handles messenger messages
type DefaultHandler struct {
	mg *Messenger
}

// ProcessMessage messages from messenger
func (h *DefaultHandler) ProcessMessage(m *messaging) {
	logger.Debugf("Received message: %s", m)
	h.mg.MarkSeen(&m.Sender)
	h.mg.SendTypingOn(&m.Sender)

	st := chatbot.NewStatement(m.Message.Text, m.Sender.ID)
	output := h.mg.Bot.GetResponse(st)

	var mids []string
	for _, msg := range output.Responses {
		if res, err := h.mg.SendMessage(msg.(utils.Message)); err == nil {
			mids = append(mids, res.MessageID)
		}
	}

	// log outgoing message
	msg := models.NewMessage(m.Sender.ID, output.Text, output.SerializeResponse(), output.Meta.String(), mids)
	msg.Save()
}

// ProcessPostback postback from messenger
func (h *DefaultHandler) ProcessPostback(p *messaging) {
	logger.Debugf("Received postback: %s", p)
	h.mg.MarkSeen(&p.Sender)
	h.mg.SendTypingOn(&p.Sender)

	st := chatbot.NewStatement(p.Postback.Title, p.Sender.ID)
	st.SetPayload(p.Postback.Payload)

	output := h.mg.Bot.GetResponse(st)

	var mids []string
	for _, msg := range output.Responses {
		if res, err := h.mg.SendMessage(msg.(utils.Message)); err == nil {
			mids = append(mids, res.MessageID)
		}
	}

	// log outgoing message
	msg := models.NewMessage(p.Sender.ID, output.Text, output.SerializeResponse(), output.Meta.String(), mids)
	msg.Save()
}

// ProcessDelivery delivery response from messenger
func (h *DefaultHandler) ProcessDelivery(d *messaging) {
	logger.Debugf("Message delivered: %s", d)

	models.MarkMessageDelivered(d.Delivery.Mids)
}
