package crawler

import (
	"fmt"
	"github.com/epigos/newsbot/utils"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	delayInterval         = 5
	responseBuffer        = 10
	topStories            = "Top stories"
	politicsCategory      = "Politics"
	worldCategory         = "World"
	sportsCategory        = "Sports"
	businessCategory      = "Business"
	lifestyleCategory     = "Lifestyle"
	entertainmentCategory = "Entertainment"
	techCategory          = "Tech"
	africaCategory        = "Africa"
)

// Crawler contains spiders to be crawled
type Crawler struct {
	Spiders       []Spider
	Logger        *utils.Logger
	resCh         chan *crawlResponse
	wg            sync.WaitGroup
	stopCh        chan bool
	ops           uint64
	crawlInterval time.Duration
}

type link struct {
	category string
	url      string
}

type links []*link

//crawlResponse response of feed request
type crawlResponse struct {
	spider   Spider
	link     *link
	response interface{}
}

// Spider an interface for spiders
type Spider interface {
	getName() string
	getLinks() links
	makeRequest(l *link)
	process(r *crawlResponse)
	setCrawler(c *Crawler)
}

func newLink(c, url string) *link {
	return &link{c, url}
}

func (l *link) String() string {
	return fmt.Sprintf("%s: %s", l.category, l.url)
}

// New creates a new crawler
func New() *Crawler {

	strInv := os.Getenv("CRAWL_INTERVAL")
	crawlInterval, err := time.ParseDuration(strInv)
	if err != nil {
		crawlInterval = time.Minute * 60
	}

	spiders := []Spider{
		newCitinews(),
		newMyjoyOnline(),
		newModernghana(),
		newGhanaweb(),
		newPulse(),
		newBBC(),
	}
	return &Crawler{
		Spiders:       spiders,
		Logger:        utils.NewLogger("crawler"),
		resCh:         make(chan *crawlResponse, responseBuffer),
		wg:            sync.WaitGroup{},
		stopCh:        make(chan bool),
		crawlInterval: crawlInterval,
	}
}

func (c *Crawler) String() string {
	return fmt.Sprintf("Crawler: %s", c.Spiders)
}

// Run starts crawling
func (c *Crawler) Run() {
	c.Logger.Info("Starting crawler")
	// reset counter
	atomic.StoreUint64(&c.ops, 0)

	for _, spider := range c.Spiders {
		c.wg.Add(1)
		spider.setCrawler(c)
		go c.Crawl(spider)
	}

	c.wg.Wait()
	c.Done()
}

// Listen listens to crawl response
func (c *Crawler) Listen() {
	for {
		select {
		case r := <-c.resCh:
			// process response
			go r.spider.process(r)

		case <-c.stopCh:
			close(c.resCh)
			c.Logger.Info("Stoping crawler")
			return
		}
	}
}

// Crawl feed links
func (c *Crawler) Crawl(s Spider) {
	links := s.getLinks()
	c.Logger.Infof("Starting spider:%s with: %v links", s.getName(), len(links))

	defer c.wg.Done()

	for _, link := range links {
		c.wg.Add(1)
		go s.makeRequest(link)
	}
}

// Stop stops crawling
func (c *Crawler) Stop() {
	c.stopCh <- true
}

// Done done crawling
func (c *Crawler) Done() {
	opsFinal := atomic.LoadUint64(&c.ops)
	c.Logger.Infof("Done crawling %d news feed", opsFinal)

	go func() {
		time.Sleep(c.crawlInterval)
		c.Run()
	}()
}

// AddResponse add crawl response
func (c *Crawler) AddResponse(cr *crawlResponse) {
	c.wg.Add(1)
	c.resCh <- cr
}
