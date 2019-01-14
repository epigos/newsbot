package models

import (
	"math/rand"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

// TopicKind kind name for topics
const TopicKind = "Topics"

// Topic is a category for article
type Topic struct {
	ID      string    `datastore:"-" json:"id"`
	Name    string    `json:"name"`
	Tags    []string  `json:"tags" datastore:",noindex"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// Key get key for article
func (m *Topic) Key() *datastore.Key {
	// if there is no Id, we want to generate an "incomplete"
	// one and let datastore determine the key/Id for us
	if m.ID == "" {
		return DS.NewKey(TopicKind)
	}

	// if Id is already set, we'll just build the Key based
	// on the one provided.
	return datastore.NameKey(TopicKind, m.ID, nil)
}

// SetID set id
func (m *Topic) SetID(key *datastore.Key) {
	m.ID = key.Name
}

func (m *Topic) String() string {
	return strings.Title(m.Name)
}

// NewTopic returns new topic
func NewTopic(n string, ts []string) *Topic {
	return &Topic{ID: strings.ToLower(n), Name: n, Tags: ts, Created: time.Now()}
}

// GetTopic get or create topic from database
func GetTopic(name string) (*Topic, error) {
	entity := Topic{ID: strings.ToLower(name)}
	err := DS.GetByKey(&entity)
	return &entity, err
}

// GetOrCreateTopic get or create topic
func GetOrCreateTopic(name string, ts []string) *Topic {
	topic, err := GetTopic(name)
	if err != nil {
		topic = NewTopic(name, ts)
		topic.Save()
		return topic
	}
	return topic
}

// Save topic
func (m *Topic) Save() {
	DS.Logger.Info("Saving topic:", m)
	DS.Save(m)
}

// GetTopicKey get or create topic from database
func GetTopicKey(name string) *datastore.Key {
	entity := Topic{ID: strings.ToLower(name)}
	return entity.Key()
}

// GetTopics get all topics
func GetTopics() ([]*Topic, error) {

	query := NewBaseQuery(TopicKind, []*Filter{})
	var topics []*Topic

	keys, err := DS.GetAll(query, &topics)
	for i, key := range keys {
		topics[i].SetID(key)
	}

	return ShuffleTopics(topics), err
}

// ShuffleTopics shuffles topics
func ShuffleTopics(topics []*Topic) []*Topic {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range topics {
		newPosition := r.Intn(len(topics) - 1)
		topics[i], topics[newPosition] = topics[newPosition], topics[i]
	}
	return topics
}
