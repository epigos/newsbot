package utils

const (
	// TopStories topic
	TopStories = "Top stories"
	// ArticleCType content type
	ArticleCType = "article"
	// GetStartedMsg message for facebook
	GetStartedMsg = "Hi, let's get you started"
	// GreetingTextMsg message to displayed at facebook page welcome screen
	GreetingTextMsg = "Welcome to News Bot {{user_first_name}}. I'll send you top stories every day, or you can ask me about a topic you want to learn more about."
	// PostBackGetStarted get started title
	PostBackGetStarted = "Get Started"
	// PostBackGetSummary get summary title
	PostBackGetSummary = "Summary"
	// PostBackShare postback button
	PostBackShare = "Share"
	// ActionNewsSearch news search action
	ActionNewsSearch = "news.search"
	// ActionNewsSearchNext news search next action
	ActionNewsSearchNext = "news.search.next"
	// ActionNewsSearchPrevious news search previous action
	ActionNewsSearchPrevious = "news.search.previous"
	// ActionNewsSearchRepeat news search repeat action
	ActionNewsSearchRepeat = "news.search.repeat"
	// ActionStop stops subscription
	ActionStop = "stop"
	// ActionReset resets subscription
	ActionReset = "reset"
	// ActionTopics list topics
	ActionTopics = "topics.lists"
	// ActionSubscribe subscribe action
	ActionSubscribe = "subscribe"
	// ActionManageAlerts manages alerts
	ActionManageAlerts = "manage.alerts"
	// SubscribeText subscribe to top stores message
	SubscribeText = "I can message you every day with top stories around the country or about a topic you're interested in."
	// ResetMsg reset message
	ResetMsg = "Hi, let's start over again."
	// SummaryScore score for reading summary
	SummaryScore = 0.01
	// ViewScore scroe for viewing article
	ViewScore = 0.02
	// NoSubscriptionText no subs text
	NoSubscriptionText = "You currently don't have any subscriptions"
	dev                = "dev"
	prod               = "prod"
	local              = "local"
)
