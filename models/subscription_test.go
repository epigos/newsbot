package models

import (
	"testing"

	"github.com/icrowley/fake"

	"github.com/stretchr/testify/assert"
)

func TestSubscription(t *testing.T) {
	assert := assert.New(t)

	bu := NewUser(fake.Characters(), fake.FirstName(), fake.LastName(), fake.DomainName(), fake.Language(), fake.Gender(), 0)
	bu.Save()

	topic := fake.Word()
	sub := NewSubscription(bu.ID, topic)
	sub.Save()
	assert.Equal(sub.String(), topic)

	subs, err := GetUserSubscriptions(bu.ID)
	assert.NoError(err)
	assert.True(len(subs) > 0)

	topics := GetUnsubscribedTopics(bu.ID, 5)
	assert.True(len(topics) > 0)
}
