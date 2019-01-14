package models

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/epigos/newsbot/utils"

	humanize "github.com/dustin/go-humanize"

	"cloud.google.com/go/datastore"
)

//ArticleKind kind name for articles
const (
	ArticleKind = "Articles"
	pageSize    = 6
)

// Assessment an Assessment provides comprehensive access to a article's metrics
type Assessment struct {
	// Automated read
	AutomatedReadability float64 `json:"automated_readability,omitempty" datastore:",noindex"`
	ColemanLiau          float64 `json:"coleman_liau,omitempty" datastore:",noindex"`
	FleschKincaid        float64 `json:"flesch_kincaid,omitempty" datastore:",noindex"`
	NumCharacters        float64 `json:"num_chars,omitempty" datastore:",noindex"`
	NumComplexWords      float64 `json:"num_complex_words,omitempty" datastore:",noindex"`
	NumParagraphs        float64 `json:"num_paragraphs,omitempty" datastore:",noindex"`
	NumPolysylWords      float64 `json:"num_polysyl_words,omitempty" datastore:",noindex"`
	NumSentences         float64 `json:"num_sentences,omitempty" datastore:",noindex"`
	NumSyllables         float64 `json:"num_syllables,omitempty" datastore:",noindex"`
	NumWords             float64 `json:"num_words,omitempty" datastore:",noindex"`
	ReadingTime          string  `json:"reading_time,omitempty" datastore:",noindex"`
}

// Article structs for crawled article
type Article struct {
	ID          string         `json:"id" datastore:"-"`
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty" datastore:",noindex"`
	Summary     []string       `json:"summary,omitempty" datastore:",noindex"`
	Link        string         `json:"link"`
	Domain      string         `json:"domain"`
	TopicKey    *datastore.Key `json:"topic"`
	Author      string         `json:"author,omitempty" datastore:",noindex"`
	Image       string         `json:"image" datastore:",noindex"`
	Tags        []string       `json:"tags,omitempty"`
	Assessment  *Assessment    `json:"assessment,omitempty" datastore:",noindex"`
	Score       float64        `json:"score"`
	Published   *time.Time     `json:"published,omitempty"`
	Created     time.Time      `json:"created"`
	Updated     time.Time      `json:"updated"`
}

// NewArticle returns article
func NewArticle(title, guid, desc, link, domain, img string, pub *time.Time, ts []string) *Article {
	return &Article{
		ID:          guid,
		Title:       title,
		Description: desc,
		Link:        link,
		Domain:      domain,
		Image:       img,
		Tags:        ts,
		Published:   pub,
		Created:     time.Now(),
	}
}

func (m *Article) String() string {
	return m.Title
}

// AddAssessment add assessment to article
func (m *Article) AddAssessment(ta *utils.TextAnalysis) {

	as := ta.Doc.Assess()

	m.Assessment = &Assessment{
		AutomatedReadability: as.AutomatedReadability,
		ColemanLiau:          as.ColemanLiau,
		FleschKincaid:        as.FleschKincaid,
		NumCharacters:        ta.Doc.NumCharacters,
		NumComplexWords:      ta.Doc.NumComplexWords,
		NumParagraphs:        ta.Doc.NumParagraphs,
		NumPolysylWords:      ta.Doc.NumPolysylWords,
		NumSentences:         ta.Doc.NumSentences,
		NumSyllables:         ta.Doc.NumSyllables,
		NumWords:             ta.Doc.NumWords,
		ReadingTime:          utils.HumanizeDuration(ta.ReadingTime()),
	}
}

// SetTopic set article topic
func (m *Article) SetTopic(name string, ts []string) {
	topic := GetOrCreateTopic(name, ts)
	m.TopicKey = topic.Key()
}

// Key get key for article
func (m *Article) Key() *datastore.Key {
	// if there is no Id, we want to generate an "incomplete"
	// one and let datastore determine the key/Id for us
	if m.ID == "" {
		return DS.NewKey(ArticleKind)
	}

	// if Id is already set, we'll just build the Key based
	// on the one provided.
	return datastore.NameKey(ArticleKind, m.ID, nil)
}

// SetID set id
func (m *Article) SetID(key *datastore.Key) {
	m.ID = key.Name
}

// GetArticleKey get article key
func GetArticleKey(id string) *datastore.Key {
	entity := Article{ID: id}
	return entity.Key()
}

// GetArticle article
func GetArticle(id string) (*Article, error) {
	entity := Article{ID: id}
	err := DS.GetByKey(&entity)
	return &entity, err
}

// Save article
func (m *Article) Save() {
	DS.Logger.Info("Saving article:", m)
	DS.Save(m)
}

// Delete article
func (m *Article) Delete() error {
	return DS.Delete(m.Key())
}

// SearchArticle search article based on params from dialogflow
func SearchArticle(params utils.Map, page int) ([]*Article, error) {
	var filters []*Filter
	var query *Query
	var articles []*Article

	if c := params.Get("category", ""); c != "" {
		name := strings.Title(c.(string))
		topic := GetTopicKey(name)
		filters = append(filters, NewFilter("TopicKey =", topic))
	}
	if kwd := params.Get("keyword", ""); kwd != "" {
		for _, word := range strings.Fields(kwd.(string)) {
			q := strings.ToLower(strings.TrimSpace(word))
			filters = append(filters, NewFilter("Tags =", q))
		}
	}
	if dt := params.Get("date-time", ""); dt != "" {
		filters = append(filters, NewFilter("Published >=", dt))
	}
	if src := params.Get("source", ""); src != "" {
		filters = append(filters, NewFilter("Domain =", src))
	}

	// send top stories if no filters available
	if len(filters) < 1 {
		topic := GetTopicKey(utils.TopStories)
		filters = append(filters, NewFilter("TopicKey =", topic))
	}

	query = NewQuery(ArticleKind, filters, pageSize, page, "-Published")

	keys, err := DS.GetAll(query, &articles)
	for i, key := range keys {
		articles[i].SetID(key)
	}
	return articles, err
}

// ToMessengerElement converts news article to messenger template
func (m *Article) ToMessengerElement(userID string) *utils.Element {
	bs := []*utils.Button{
		utils.NewPostbackButton(utils.PostBackGetSummary, m.ID),
		utils.NewWebURLButton("Read More", m.GetMessengerLink(userID)),
		utils.NewShareButton(),
	}
	return utils.NewElement(m.Title, m.GetSubText(), m.Link, m.Image, bs)
}

// GetMessengerLink generate messenger link
func (m *Article) GetMessengerLink(userID string) string {
	host := os.Getenv("APP_HOST")
	return fmt.Sprintf("%s/ns/%s/%s", host, m.Key().Encode(), userID)
}

// GetSubText generate subtitle for messenger element
func (m *Article) GetSubText() string {
	return fmt.Sprintf("%s • %s • %s read", m.TopicKey.Name, humanize.Time(*m.Published), m.Assessment.ReadingTime)
}
