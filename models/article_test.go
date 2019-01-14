package models

import (
	"github.com/epigos/newsbot/utils"
	"testing"
	"time"

	"github.com/icrowley/fake"

	"github.com/stretchr/testify/assert"
)

func TestArticle(t *testing.T) {
	assert := assert.New(t)

	pub := time.Now()
	ts := []string{fake.Word(), fake.Word()}
	link := fake.DomainName()
	topic := fake.Word()

	title := fake.SentencesN(1)
	nw := NewArticle(title, link, fake.SentencesN(2), link, link, link, &pub, ts)

	ta := utils.NewTextAnalysis(fake.SentencesN(10), fake.SentencesN(10))
	nw.AddAssessment(ta)
	nw.SetTopic(topic, ts)
	nw.Save()

	assert.Equal(nw.ID, link)
	assert.Equal(nw.Title, title)
	assert.Equal(nw.Published, &pub)
	assert.Equal(nw.Tags, ts)
	assert.Equal(nw.TopicKey.Name, topic)

	el := nw.ToMessengerElement("id")
	assert.Equal(el.Title, nw.Title)
	assert.Equal(len(el.Buttons), 3)

	article, err := GetArticle(link)
	assert.NoError(err)
	assert.Equal(article.Title, title)

	l := nw.GetMessengerLink("id")
	assert.Contains(l, "ns")

	nw.Title = fake.SentencesN(1)
	nw.Save()
	assert.NotEqual(nw.Title, title)

	err = nw.Delete()
	assert.NoError(err)
}

func TestSearchArticle(t *testing.T) {
	assert := assert.New(t)

	pub := time.Now()
	ts := []string{fake.Word(), fake.Word()}
	link := fake.DomainName()
	topic := fake.Word()

	title := fake.SentencesN(1)
	nw := NewArticle(title, link, fake.SentencesN(2), link, link, link, &pub, ts)
	ta := utils.NewTextAnalysis(fake.SentencesN(10), fake.SentencesN(10))
	nw.AddAssessment(ta)
	nw.SetTopic(topic, ts)
	nw.Save()

	m := utils.Map{
		"keyword":   ts[0],
		"date-time": pub,
		"category":  topic,
		"source":    link,
	}

	articles, err := SearchArticle(m, 1)
	assert.NoError(err)
	assert.NotEmpty(articles)

	m = m.Remove("category")
	articles, err = SearchArticle(m, 1)
	assert.NoError(err)
	assert.NotEmpty(articles)
}
