package models

import (
	"testing"

	"github.com/icrowley/fake"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	assert := assert.New(t)

	fn, ln := fake.FirstName(), fake.LastName()
	id := fake.Characters()
	user := NewUser(id, fn, ln, fake.DomainName(), fake.Language(), fake.Gender(), 0)
	user.Save()
	assert.Equal(user.String(), fn+" "+ln)

	nUser, err := GetUser(user.ID)
	assert.NoError(err)
	assert.Equal(user.FirstName, nUser.FirstName)

	nUser.FirstName = "firstName"
	nUser.Save()
	assert.NotEqual(nUser.FirstName, fn)
	assert.Equal(nUser.FirstName, "firstName")
}

func TestUserAction(t *testing.T) {
	assert := assert.New(t)

	user := NewUser(fake.Characters(), fake.FirstName(), fake.LastName(), fake.DomainName(), fake.Language(), fake.Gender(), 0)
	user.Save()

	userAction := NewUserAction(user.Key().Name, user.Key(), "view")
	userAction.Save()
	assert.Equal(userAction.UserKey, user.Key())
}
