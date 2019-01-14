package messenger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSenderAction(t *testing.T) {
	assert := assert.New(t)
	seen := NewSenderAction(string(rid), ActionMarkSeen)
	assert.Equal(string(rid), seen.Recipient.ID)
	assert.Equal(ActionMarkSeen, seen.Action)
}
