package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTextMessage(t *testing.T) {
	assert := assert.New(t)

	m := NewTextMessage("id", "Hi, there!")

	assert.Equal(m.Message.Text, "Hi, there!")
	assert.Equal(m.Recipient.ID, "id")
}

func TestGenericMessage(t *testing.T) {
	assert := assert.New(t)

	m := NewGenericMessage("id")
	assert.Equal(m.Recipient.ID, "id")
	el := NewElement("title", "subtitle", "itemUrl", "imageURL", []*Button{})
	m.AddElement(el)
	assert.Contains(m.Message.Attachment.Payload.Elements, el)
	m = NewGenericMessage("id")
	m.AddNewElement(el.Title, el.Subtitle, el.ItemURL, el.ImageURL, el.Buttons)
	assert.Contains(m.Message.Attachment.Payload.Elements, el)
}

func TestNewWebURLButton(t *testing.T) {
	assert := assert.New(t)
	m := NewWebURLButton("example", "http://example.com")
	assert.Equal(m.Title, "example")
	assert.Equal(m.URL, "http://example.com")
	assert.Equal(m.Type, ButtonTypeWebURL)
}

func TestNewPostbackButton(t *testing.T) {
	assert := assert.New(t)
	m := NewPostbackButton("example", "payload")
	assert.Equal(m.Title, "example")
	assert.Equal(m.Payload, "payload")
	assert.Equal(m.Type, ButtonTypePostback)
}

func TestElement(t *testing.T) {
	assert := assert.New(t)
	m := NewElement("title", "subtitle", "itemUrl", "imageURL", []*Button{})
	assert.Equal(m.Title, "title")
	assert.Equal(m.Subtitle, "subtitle")
	assert.Equal(m.ImageURL, "imageURL")
	assert.Equal(m.ItemURL, "itemUrl")
	assert.Equal(m.Buttons, []*Button{})

	wb := NewWebURLButton("example", "http://example.com")
	pb := NewPostbackButton("example", "payload")
	m.AddWebURLButton(wb.Title, wb.URL)
	assert.Contains(m.Buttons, wb)
	m.AddPostbackButton(pb.Title, pb.Payload)
	assert.Contains(m.Buttons, pb)
}
