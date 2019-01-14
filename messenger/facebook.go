package messenger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/web"

	"github.com/mitchellh/mapstructure"
)

// Recipient represents facebook recipient
type Recipient struct {
	ID      string       `json:"id"`
	Profile *models.User `json:"-" bson:"-"`
}

func (r Recipient) String() string {
	if r.Profile != nil {
		return r.Profile.String()
	}
	return r.ID
}

type messaging struct {
	Recipient Recipient         `json:"recipient"`
	Sender    Recipient         `json:"sender"`
	Timestamp int               `json:"timestamp"`
	Message   *FacebookMessage  `json:"message,omitempty"`
	Delivery  *FacebookDelivery `json:"delivery"`
	Postback  *FacebookPostback `json:"postback"`
}

func (m *messaging) String() string {
	if m.Message != nil {
		return fmt.Sprintf("From: %s, Text: %s", m.Sender, m.Message.Text)
	} else if m.Postback != nil {
		return fmt.Sprintf("From: %s, Title: %s", m.Sender, m.Postback.Title)
	} else if m.Delivery != nil {
		return strings.Join(m.Delivery.Mids, ", ")
	}
	return ""
}

// Mid facebook message id
func (m *messaging) Mid() string {
	if m.Message != nil {
		return m.Message.Mid
	}
	return ""
}

type facebookEntry struct {
	ID        string      `json:"id"`
	Messaging []messaging `json:"messaging"`
	Time      int         `json:"time"`
}

// FacebookRequest received from Facebook server on webhook, contains messages, delivery reports and/or postbacks
type FacebookRequest struct {
	Entry  []facebookEntry `json:"entry"`
	Object string          `json:"object"`
}

// FacebookMessage struct for text messaged received from facebook server as part of FacebookRequest struct
type FacebookMessage struct {
	Mid        string               `json:"mid"`
	Seq        int                  `json:"seq"`
	Text       string               `json:"text"`
	Attachment []facebookAttachment `json:"attachments"`
}

type facebookAttachment struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// FacebookDelivery struct for delivery reports received from Facebook server as part of FacebookRequest struct
type FacebookDelivery struct {
	Mids      []string `json:"mids"`
	Seq       int      `json:"seq"`
	Watermark int      `json:"watermark"`
}

// FacebookPostback struct for postbacks received from Facebook server  as part of FacebookRequest struct
type FacebookPostback struct {
	Title   string `json:"title"`
	Payload string `json:"payload"`
}

// rawFBResponse received from Facebook server after sending the message
// if Error is null we copy this into FacebookResponse object
type rawFBResponse struct {
	MessageID   string         `json:"message_id"`
	RecipientID string         `json:"recipient_id"`
	Error       *FacebookError `json:"error"`
}

// FacebookResponse received from Facebook server after sending the message
type FacebookResponse struct {
	MessageID   string `json:"message_id"`
	RecipientID string `json:"recipient_id"`
}

func (r *FacebookResponse) String() string {
	return fmt.Sprintf("FB response: %s %v", r.MessageID, r.RecipientID)
}

// FacebookError received form Facebook server if sending messages failed
type FacebookError struct {
	Code      int    `json:"code"`
	FbtraceID string `json:"fbtrace_id"`
	Message   string `json:"message"`
	Type      string `json:"type"`
}

// Error returns Go error object constructed from FacebookError data
func (err *FacebookError) Error() error {
	return fmt.Errorf("FB Error: Type %s: %s; FB trace ID: %s", err.Type, err.Message, err.FbtraceID)
}

// DecodeRequest decodes http request from FB messagner to FacebookRequest struct
// DecodeRequest will close the Body reader
// Usually you don't have to use DecodeRequest if you setup events for specific types
func (mg *Messenger) DecodeRequest(ctx *web.Context) (*FacebookRequest, error) {
	var fbRq FacebookRequest
	err := mapstructure.Decode(ctx.PostValues(), &fbRq)
	return &fbRq, err
}

// decodeResponse decodes Facebook response after sending message, usually contains MessageID or Error
func (mg *Messenger) decodeResponse(r *http.Response) (*FacebookResponse, error) {
	defer r.Body.Close()
	var fbResp rawFBResponse
	err := json.NewDecoder(r.Body).Decode(&fbResp)
	if err != nil {
		logger.Error(err)
		return &FacebookResponse{}, err
	}

	if fbResp.Error != nil {
		logger.Error(fbResp.Error.Error())
		return &FacebookResponse{}, fbResp.Error.Error()
	}

	re := &FacebookResponse{
		MessageID:   fbResp.MessageID,
		RecipientID: fbResp.RecipientID,
	}
	logger.Debugf("%s", re)

	return re, nil
}
