package chatbot

import (
	"log"
	"testing"

	"github.com/epigos/newsbot/models"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"

	"github.com/joho/godotenv"
)

var (
	user *models.User
	ch   *Chatbot
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	n := "Test"

	assert.Equal(ch.Name, n, "they should be equal")

	if _, ok := ch.Logic.(LogicAdapter); ok == false {
		t.Error("Chatbot Logic adapter does not implement LogicAdapter")
	}
}

func TestGetResponse(t *testing.T) {
	assert := assert.New(t)
	n := "Hi"

	st := NewStatement(n, user.ID)
	res := ch.GetResponse(st)
	assert.Equal(res.Text, n)

	st = NewStatement("Get Started", user.ID)
	res = ch.GetResponse(st)
	assert.Len(res.Responses, 3)
}

func TestMain(m *testing.M) {
	// load env variables
	err := godotenv.Load("../env/test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	models.Connect()
	user = models.NewUser(fake.Characters(), fake.FirstName(), fake.LastName(), fake.DomainName(), fake.Language(), fake.Gender(), 0)

	ch = New("Test")

	m.Run()
	models.Close()
}
