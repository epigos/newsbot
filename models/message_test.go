package models

import (
	"github.com/epigos/newsbot/utils"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	assert := assert.New(t)

	fn, ln := fake.FirstName(), fake.LastName()
	user := NewUser(fake.CharactersN(10), fn, ln, fake.DomainName(), fake.Language(), fake.Gender(), 0)
	user.Save()

	p := utils.Map{"Text": "Hi"}
	text := fake.Sentence()
	mids := []string{fake.Characters(), fake.Characters()}
	om := NewMessage(user.ID, text, p.String(), p.String(), mids)
	om.Save()
	assert.Equal(user.Key(), om.User)
	assert.Equal(om.Text, text)

	err := MarkMessageDelivered(mids)
	assert.NoError(err)

	nm, err := GetMessage(om.ID)
	assert.NoError(err)
	assert.NotNil(nm.DeliveryTime)
}
