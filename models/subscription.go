package models

import (
	"fmt"
	"github.com/epigos/newsbot/utils"
	"time"

	"cloud.google.com/go/datastore"
)

//SubscriptionKind kind name for subscriptions
const SubscriptionKind = "Subscriptions"

// Subscription is a model for content subscriptions
type Subscription struct {
	ID      string         `datastore:"-" json:"id"`
	User    *datastore.Key `json:"user_id"`
	Topic   *datastore.Key `json:"topic"`
	Created time.Time      `json:"created"`
	Updated time.Time      `json:"updated"`
}

// Key get key for article
func (m *Subscription) Key() *datastore.Key {
	// if there is no Id, we want to generate an "incomplete"
	// one and let datastore determine the key/Id for us
	if m.ID == "" {
		return DS.NewKey(SubscriptionKind)
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
func (m *Subscription) SetID(key *datastore.Key) {
	m.ID = key.Encode()
}

// NewSubscription returns new subscription model
func NewSubscription(uid string, topic string) *Subscription {
	ukey := GetUserKey(uid)
	tkey := GetTopicKey(topic)
	return &Subscription{User: ukey, Topic: tkey}
}

func (m *Subscription) String() string {
	return m.Topic.Name
}

// Description description of subscription
func (m *Subscription) Description() string {
	return fmt.Sprintf("You'll receive %s news throughout the day", m.Topic.Name)
}

// StopButton get messenger stop button
func (m *Subscription) StopButton() []*utils.Button {
	title := fmt.Sprintf("Stop %s", m.Topic.Name)
	return []*utils.Button{utils.NewPostbackButton(title, title)}
}

// Save saves subscription
func (m *Subscription) Save() {
	DS.Logger.Info("Saving subscription:", m)
	DS.Save(m)
}

// Delete deletes subscription
func (m *Subscription) Delete() {
	DS.Logger.Info("Deleting subscription:", m)
	DS.Delete(m.Key())
}

// GetUserSubscriptions get user subscriptions
func GetUserSubscriptions(uid string) ([]*Subscription, error) {
	var fs []*Filter

	ukey := GetUserKey(uid)
	fs = []*Filter{NewFilter("User =", ukey)}

	query := NewQuery(SubscriptionKind, fs, 0, 0)
	var subs []*Subscription

	keys, err := DS.GetAll(query, &subs)
	for i, key := range keys {
		subs[i].SetID(key)
	}
	return subs, err

}

// GetUnsubscribedTopics get unsubscribed topics
func GetUnsubscribedTopics(uid string, limit int) []*Topic {
	ts, _ := GetTopics()
	subs, _ := GetUserSubscriptions(uid)

	if len(subs) < 1 {
		return ts[:limit]
	}

	var topics []*Topic

	seen := map[string]bool{}

	for _, topic := range ts {
		for _, sub := range subs {
			if _, ok := seen[topic.ID]; !ok {
				if sub.Topic != topic.Key() && len(topics) < limit {
					topics = append(topics, topic)
					seen[topic.ID] = true
				}
			}
		}
	}
	return topics
}
