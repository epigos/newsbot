package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/epigos/newsbot/chatbot"
	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/utils"
	"github.com/epigos/newsbot/web"
	"os"
	"time"
)

const (
	bufSize      = 1024
	apiURL       = "https://graph.facebook.com/v2.6"
	messagesPath = "messages"
	profilePath  = "profile"
)

// TestURL to mock FB server, used for testing
var (
	TestURL   = ""
	endPoints = map[string]string{
		messagesPath: "me/messages",
		profilePath:  "me/messenger_profile",
	}
	logger = utils.NewLogger("messenger")
)

// Messenger struct
type Messenger struct {
	AccessToken string
	VerifyToken string
	PageID      string
	messageCh   chan *messaging
	deliveryCh  chan *messaging
	postbackCh  chan *messaging
	Logger      *utils.Logger
	Bot         *chatbot.Chatbot
	// message handler
	Handler MessageHandler
	PushCh  chan bool
}

// New creates new messenger instance
func New(b *chatbot.Chatbot) *Messenger {
	messageCh := make(chan *messaging, bufSize)
	deliveryCh := make(chan *messaging, bufSize)
	postbackCh := make(chan *messaging, bufSize)

	m := &Messenger{
		AccessToken: os.Getenv("FACEBOOK_ACCESS_TOKEN"),
		VerifyToken: os.Getenv("FACEBOOK_SECRET_TOKEN"),
		PageID:      os.Getenv("FACEBOOK_PAGE_ID"),
		// messageCh channel for events when message from Facebook is received
		messageCh: messageCh,
		// deliveryCh channel for events when delivery report from Facebook received
		deliveryCh: deliveryCh,
		// postbackCh channel for events when postback received from Facebook
		postbackCh: postbackCh,
		Logger:     logger,
		Bot:        b,
		PushCh:     make(chan bool),
	}
	m.Handler = &DefaultHandler{m}
	return m
}

func (mg *Messenger) buildURL(path string) string {
	p, ok := endPoints[path]
	if !ok {
		p = path
	}
	var url = "%s/%s?access_token=%s"

	if TestURL != "" {
		url = fmt.Sprintf(url, TestURL, p, mg.AccessToken)
	} else {
		url = fmt.Sprintf(url, apiURL, p, mg.AccessToken)
	}

	return url
}

// VerifyWebhook verifies your webhook by checking VerifyToken and sending challange back to Facebook
func (mg *Messenger) VerifyWebhook(ctx *web.Context) *web.HTTPError {
	// Facebook sends this query for verifying webhooks
	// hub.mode=subscribe&hub.challenge=1085525140&hub.verify_token=moj_token
	if ctx.FormValue("hub.mode") == "subscribe" {
		if ctx.FormValue("hub.verify_token") == mg.VerifyToken {
			return ctx.WriteString(ctx.FormValue("hub.challenge"))
		}
	}
	return nil
}

func (mg *Messenger) processEntry(fs []facebookEntry) {

	for _, entry := range fs {
		for _, msg := range entry.Messaging {
			// get sender profile
			msg.Sender.Profile = mg.GetSenderProfile(msg.Sender.ID)
			switch {
			case msg.Message != nil:
				mg.messageCh <- &msg
			case msg.Delivery != nil:
				mg.deliveryCh <- &msg
			case msg.Postback != nil:
				mg.postbackCh <- &msg
			}
		}
	}
}

// ServeHTTP is HTTP handler for Messenger so it could be directly used as http.Handler
func (mg *Messenger) ServeHTTP(ctx *web.Context) *web.HTTPError {
	er := mg.VerifyWebhook(ctx) // verify webhook if needed
	if er != nil {
		return er
	}
	fbRq, err := mg.DecodeRequest(ctx) // get FacebookRequest object

	if err != nil {
		e := fmt.Sprintf("Facebook request error:%v", err)
		logger.Info(e)
		return ctx.BadRequest(e)
	}
	mg.processEntry(fbRq.Entry)

	return ctx.WriteString("Message received")
}

// Listen messenger channel events
func (mg *Messenger) Listen() {
	logger.Info("Messenger channels opened")
	for {
		select {
		case m := <-mg.messageCh:
			go mg.Handler.ProcessMessage(m)
		case d := <-mg.deliveryCh:
			go mg.Handler.ProcessDelivery(d)
		case p := <-mg.postbackCh:
			go mg.Handler.ProcessPostback(p)
		case <-mg.PushCh:
			go mg.PushMessages()
		}
	}
}

// makeFbRequest makes request to facebook
func (mg *Messenger) makeFbRequest(path, method string, body interface{}) (*http.Response, error) {
	url := mg.buildURL(path)
	var s []byte
	if body != nil {
		s, _ = json.Marshal(body)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(s))
	req.Header.Set("Content-Type", "application/json")

	logger.Debugf("Making FB request to %s; method: %s; params: %s", url, method, string(s))
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

// GetSenderProfile returns facebook profile
func (mg *Messenger) GetSenderProfile(senderID string) *models.User {

	user, err := models.GetUser(senderID)
	if err == nil {
		return user
	}
	logger.Info(err)

	r, err := mg.makeFbRequest(senderID, "GET", nil)

	if err != nil {
		logger.Error(err)
		return user
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		logger.Error(err)
	}
	// save user
	user.Created = time.Now()
	user.Save()
	if err != nil {
		logger.Error(err)
	}
	return user
}

// PushMessages push new crawled items to users
func (mg *Messenger) PushMessages() {
	logger.Info("Crawl done,", "push messages")
}

// SendTextMessage sends text messate to receiverID
// it is shorthand instead of crating new text message and then sending it
func (mg *Messenger) SendTextMessage(receiverID string, text string) (*FacebookResponse, error) {
	m := utils.NewTextMessage(receiverID, text)
	return mg.SendMessage(m)
}

// SendMessage sends chat message
func (mg *Messenger) SendMessage(m utils.Message) (*FacebookResponse, error) {
	resp, err := mg.makeFbRequest(messagesPath, "POST", m)

	if err != nil {
		return &FacebookResponse{}, err
	}
	return mg.decodeResponse(resp)
}
