package crawler

import (
	"fmt"
	"log"
	"testing"

	"github.com/epigos/newsbot/models"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// load env variables
	err := godotenv.Load("../env/test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	models.Connect()
	m.Run()
	models.Close()
}

func TestCrawler(t *testing.T) {
	assert := assert.New(t)

	testSpider := newFeedSpider(
		"citinewsroom",
		"citinewsroom.com",
		newLink(topStories, "https://citinewsroom.com/ghana-news/top-stories/feed/"),
		newLink(worldCategory, "http://example.com/feed/"),
	)
	testSpider.Config.Set("UseMetaDesc", true)
	testSpider.Config.Set("BodySelector", ".entry-content p")

	c := New()
	c.Spiders = []Spider{testSpider}

	// ch := make(chan bool)
	go c.Run()
	c.Stop()

	assert.Contains(c.String(), testSpider.getName())
	// done := <-ch
	// assert.True(done)

	// articles, err := models.SearchArticle(utils.Map{}, 5, 0)
	// assert.NoError(err)
	// assert.NotEmpty(articles)
	// for _, article := range articles {
	// 	assert.Equal(article.Domain, "citinewsroom.com")
	// }
}
func TestLink(t *testing.T) {
	assert := assert.New(t)

	ln := newLink(topStories, "http://example.com")
	assert.Equal(ln.url, "http://example.com")
	assert.Equal(ln.category, topStories)
	assert.Equal(ln.String(), fmt.Sprintf("%s: %s", ln.category, ln.url))
}

func TestSpider(t *testing.T) {
	assert := assert.New(t)

	sm := newMyjoyOnline()
	assert.Equal(sm.getName, "myjoyonline")

	cn := newCitinews()
	assert.Equal(cn.getName(), "citinewsroom")

	bb := newBBC()
	assert.Equal(bb.getName(), "bbc")

	gw := newGhanaweb()
	assert.Equal(gw.getName(), "ghanaweb")

	mg := newModernghana()
	assert.Equal(mg.getName(), "modernghana")

	pl := newPulse()
	assert.Equal(pl.getName(), "pulse")
}
