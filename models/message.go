package models

import (
	"time"

	"cloud.google.com/go/datastore"
)

//MessageKind kind name for messages
const MessageKind = "Messages"

// Message recieved from facebook
type Message struct {
	ID           string         `datastore:"-" json:"id"`
	User         *datastore.Key `json:"user_id"`
	Text         string         `json:"text"`
	MID          []string       `json:"mids"`
	Response     string         `datastore:",noindex"  json:"response"`
	Meta         string         `datastore:",noindex"  json:"meta"`
	DeliveryTime *time.Time     `json:"delivery_time"`
	Created      time.Time      `json:"created"`
	Updated      time.Time      `json:"updated"`
}

// Key get key for article
func (m *Message) Key() *datastore.Key {
	// if there is no Id, we want to generate an "incomplete"
	// one and let datastore determine the key/Id for us
	if m.ID == "" {
		return DS.NewKey(MessageKind)
	}

	// if Id is already set, we'll just build the Key based
	// on the one provided.
	key, err := datastore.DecodeKey(m.ID)
	if err != nil {
		DS.Logger.Error("Key not found:", err)
	}
	return key
}

// SetID set id
func (m *Message) SetID(key *datastore.Key) {
	m.ID = key.Encode()
}

func (m *Message) String() string {
	return m.Text
}

// NewMessage creates a new model
func NewMessage(rid, text, res, meta string, mids []string) *Message {
	ukey := GetUserKey(rid)

	msg := Message{User: ukey, Text: text, MID: mids, Response: res, Meta: meta}
	return &msg
}

// GetMessage get message by id
func GetMessage(id string) (*Message, error) {
	entity := Message{ID: id}
	err := DS.GetByKey(&entity)
	return &entity, err
}

// Save messages
func (m *Message) Save() {
	DS.Logger.Info("Saving message:", m)
	DS.Save(m)
}

// MarkMessageDelivered mark outgoing messages as delivered
func MarkMessageDelivered(mids []string) error {
	var fs []*Filter

	for _, mid := range mids {
		fs = append(fs, NewFilter("MID =", mid))
	}

	query := NewQuery(MessageKind, fs, 0, 0)
	var messages []*Message

	keys, err := DS.GetAll(query, &messages)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, msg := range messages {
		msg.DeliveryTime = &now
	}
	_, err = DS.Client.PutMulti(DS.Context, keys, messages)
	return err
}
