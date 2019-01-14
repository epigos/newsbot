package crawler

import (
	"sync/atomic"
	"time"

	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/utils"

	"github.com/mmcdole/gofeed"
)

type feedSpider struct {
	Name    string
	Domain  string
	Links   links
	Config  utils.Map
	crawler *Crawler
	parser  *gofeed.Parser
}

func newFeedSpider(name, domain string, links ...*link) *feedSpider {
	return &feedSpider{
		Name:   name,
		Domain: domain,
		Config: utils.Map{},
		parser: gofeed.NewParser(),
		Links:  links,
	}
}

func (s *feedSpider) String() string {
	return s.getName()
}

func (s *feedSpider) setCrawler(c *Crawler) {
	s.crawler = c
}

func (s *feedSpider) getLinks() links {
	return s.Links
}

func (s *feedSpider) getName() string {
	return s.Name
}

func (s *feedSpider) makeRequest(l *link) {
	defer s.crawler.wg.Done()

	feed, err := s.parser.ParseURL(l.url)

	if err != nil {
		s.crawler.Logger.Debugf("%s might be down!", err)
		return
	}
	s.crawler.Logger.Debugf("%s is up!", l)
	s.crawler.AddResponse(&crawlResponse{s, l, feed})
}

// Process process items from crawl
func (s *feedSpider) process(r *crawlResponse) {
	defer s.crawler.wg.Done()

	feed := r.response.(*gofeed.Feed)
	s.crawler.Logger.Infof("Found %v items at %s", len(feed.Items), r.link.url)

	for _, i := range feed.Items {
		// delay to avoid ddos on news sites
		time.Sleep(time.Second * delayInterval)

		doc, err := utils.LinkToDoc(i.Link)

		if err != nil {
			s.crawler.Logger.Debug(err)
			continue
		}

		meta, err := utils.ExtractMetaTags(doc, "og:")
		if err != nil {
			s.crawler.Logger.Debug(err)
			continue
		}

		img := meta.Get("image", nil)
		if img == nil {
			s.crawler.Logger.Debug("Image not found: ", meta)
			continue
		}

		desc := i.Description
		useMetaDesc := s.Config.Get("UseMetaDesc", nil)
		if useMetaDesc == true {
			de := meta.Get("description", nil)
			if de == nil {
				s.crawler.Logger.Debug("description not found: ", meta)
				continue
			}
			desc = de.(string)
		}

		sel := s.Config.Get("BodySelector", nil)
		body := doc.Find(sel.(string)).Text()
		if body == "" {
			continue
		}
		ta := utils.NewTextAnalysis(body, desc)

		article := models.NewArticle(i.Title, i.GUID, desc, i.Link, s.Domain, img.(string), i.PublishedParsed, ta.Tags())
		article.SetTopic(r.link.category, []string{})

		if i.Author != nil {
			article.Author = i.Author.Name
		}

		article.Summary = ta.Sentences(3)
		article.AddAssessment(ta)

		_, er := models.GetArticle(article.ID)
		if er != nil {
			// found new article
			atomic.AddUint64(&s.crawler.ops, 1)
		}
		article.Save()
	}
}

// newBBC creates new feed spider for bbc
func newBBC() *feedSpider {
	sp := newFeedSpider(
		"bbc",
		"bbc.com",
		newLink(africaCategory, "http://feeds.bbci.co.uk/news/world/africa/rss.xml"),
		newLink(worldCategory, "http://feeds.bbci.co.uk/news/world/rss.xml"),
	)
	sp.Config.Set("UseMetaDesc", false)
	sp.Config.Set("BodySelector", ".story-body__inner p")

	return sp
}

// newCitinews creates new feed spider for citi news
func newCitinews() *feedSpider {
	sp := newFeedSpider(
		"citinewsroom",
		"citinewsroom.com",
		newLink(topStories, "https://citinewsroom.com/ghana-news/top-stories/feed/"),
		newLink(politicsCategory, "https://citinewsroom.com/ghana-news/politics/feed/"),
		newLink(sportsCategory, "https://citinewsroom.com/ghana-news/sports/feed/"),
		newLink(businessCategory, "https://citinewsroom.com/ghana-news/business/feed/"),
		newLink(entertainmentCategory, "https://citinewsroom.com/ghana-news/showbiz/feed/"),
	)
	sp.Config.Set("UseMetaDesc", true)
	sp.Config.Set("BodySelector", ".entry-content p")

	return sp
}

// newGhanaweb creates new feed spider for ghanaweb
func newGhanaweb() *feedSpider {
	sp := newFeedSpider(
		"ghanaweb",
		"ghanaweb.com",
		newLink(topStories, "https://cdn.ghanaweb.com/feed/newsfeed.xml"),
		newLink(sportsCategory, "https://cdn.ghanaweb.com/feed/soccerfeed.xml"),
		newLink(sportsCategory, "https://cdn.ghanaweb.com/feed/other_sportsfeed.xml"),
		newLink(entertainmentCategory, "https://cdn.ghanaweb.com/feed/entertainmentfeed.xml"),
	)
	sp.Config.Set("UseMetaDesc", false)
	sp.Config.Set("BodySelector", "#medsection1 > p:nth-child(8)")

	return sp
}

// newModernghana creates new feed spider for modernghana
func newModernghana() *feedSpider {
	sp := newFeedSpider(
		"modernghana",
		"modernghana.com",
		newLink(topStories, "https://rss.modernghana.com/news.xml?cat_id=1&group_id=1"),
		newLink(politicsCategory, "https://rss.modernghana.com/news.xml?cat_id=1&group_id=5"),
		newLink(sportsCategory, "https://rss.modernghana.com/news.xml?cat_id=2"),
		newLink(businessCategory, "https://rss.modernghana.com/news.xml?cat_id=1&group_id=6"),
		newLink(entertainmentCategory, "https://rss.modernghana.com/news.xml?cat_id=3"),
		newLink(worldCategory, "https://rss.modernghana.com/news.xml?cat_id=1&group_id=8"),
		newLink(africaCategory, "https://rss.modernghana.com/news.xml?cat_id=1&group_id=2"),
	)
	sp.Config.Set("UseMetaDesc", false)
	sp.Config.Set("BodySelector", ".blog-content p")

	return sp
}

// newPulse creates new feed spider for pulse.com.gh
func newPulse() *feedSpider {
	sp := newFeedSpider(
		"pulse",
		"pulse.com.gh",
		newLink(topStories, "http://www.pulse.com.gh/rss"),
	)

	sp.Config.Set("UseMetaDesc", false)
	sp.Config.Set("BodySelector", ".article_text p")

	return sp
}

// newMyjoyOnline creates new feed spider for MyjoyOnline
func newMyjoyOnline() *feedSpider {
	sp := newFeedSpider(
		"myjoyonline",
		"myjoyonline.com",
		newLink(topStories, "https://www.myjoyonline.com/pages/rss/site_edition.xml"),
		newLink(politicsCategory, "https://www.myjoyonline.com/pages/rss/site_politics.xml"),
		newLink(worldCategory, "https://www.myjoyonline.com/pages/rss/site_world.xml"),
		newLink(sportsCategory, "https://www.myjoyonline.com/pages/rss/site_sports.xml"),
		newLink(businessCategory, "https://www.myjoyonline.com/pages/rss/site_business.xml"),
		newLink(lifestyleCategory, "https://www.myjoyonline.com/pages/rss/site_lifestyle.xml"),
		newLink(entertainmentCategory, "https://www.myjoyonline.com/pages/rss/site_entertainment.xml"),
		newLink(techCategory, "https://www.myjoyonline.com/pages/rss/site_technology.xml"),
	)

	sp.Config.Set("UseMetaDesc", false)
	sp.Config.Set("BodySelector", ".article-text p")

	return sp
}
