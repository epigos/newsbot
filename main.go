package main

import (
	"flag"
	"os"

	"github.com/epigos/newsbot/chatbot"
	"github.com/epigos/newsbot/crawler"
	"github.com/epigos/newsbot/utils"

	"github.com/epigos/newsbot/messenger"
	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/web"

	"github.com/joho/godotenv"
	rollbar "github.com/rollbar/rollbar-go"
)

func main() {

	var host = flag.String("host", "0.0.0.0:5050", "host and port to run")
	var setupFbPage = flag.Bool("setup-fb-page", false, "setup facebook get started and greetings screen")
	var crawlerMode = flag.Bool("crawler", false, "start background crawler")
	flag.Parse()

	setupRollbar()
	models.Connect()
	// worker
	if *crawlerMode == true {
		cr := crawler.New()
		go cr.Run()
		cr.Listen()
	}
	// chatbot
	ch := chatbot.New("Newsbot")
	// messenger bot
	messenger := messenger.New(ch)
	// listens messenger channel events
	go messenger.Listen()
	// setup facebook screen page
	if *setupFbPage == true {
		go messenger.SetupPage()
	}
	// web server
	s := web.New(*host)
	s.Handle("/facebook", messenger.ServeHTTP, "GET", "POST")
	s.Post("/_test/bot", ch.TestHandler)
	// start server
	s.Run()
}

func setupRollbar() {
	rollbar.SetToken(os.Getenv("ROLLBAR_TOKEN"))
	rollbar.SetEnvironment(utils.GetEnvironment()) // defaults to "development"
	rollbar.SetServerHost(os.Getenv("HOST_NAME"))  // optional override; defaults to hostname
	rollbar.SetCodeVersion(utils.GetVersion())
}

func init() {
	logger := utils.NewLogger("main")
	// load env variables
	err := godotenv.Load("env/local.env")
	if err != nil {
		logger.Warn("Error loading .env file;", "using env from environment")
	}
}
