package models

import (
	"testing"

	"github.com/icrowley/fake"

	"github.com/stretchr/testify/assert"
)

func TestTopic(t *testing.T) {
	assert := assert.New(t)
	n := fake.Word()
	ts := []string{fake.Word(), fake.Word()}

	topic := NewTopic(n, ts)
	topic.Save()
	assert.Equal(topic.Name, n)
	assert.False(topic.Key().Incomplete())

	topic, err := GetTopic(n)
	assert.NoError(err)
	assert.Equal(topic.Name, n)

	nw := fake.Word()
	topic = GetOrCreateTopic(nw, ts)
	assert.NotEqual(topic.Name, n)
	assert.Equal(topic.ID, nw)

	topics, err := GetTopics()
	assert.NoError(err)
	assert.True(len(topics) > 0)
}
