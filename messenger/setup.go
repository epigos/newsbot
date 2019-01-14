package messenger

import (
	"github.com/epigos/newsbot/utils"
)

const (
	webURLType   = MenuType("web_url")
	postBackType = MenuType("postback")
	nestedType   = MenuType("nested")
)

// GetStarted Sets the Get Started button postback
type GetStarted struct {
	Content getStartedContent `json:"get_started"`
}

type getStartedContent struct {
	Payload string `json:"payload"`
}

// GreetingText Set the greeting text
type GreetingText struct {
	Greeting greeting `json:"greeting"`
}

type greeting struct {
	Locale string `json:"locale"`
	Text   string `json:"text"`
}

// MenuType persistent menu types
type MenuType string

// PersistentMenu persistent menu allows you to have an always-on
// user interface element inside Messenger conversations
type PersistentMenu struct {
	Menus []persistentMenu `json:"persistent_menu"`
}

type persistentMenu struct {
	Locale                string    `json:"locale"`
	ComposerInputDisabled bool      `json:"composer_input_disabled"`
	Actions               []*action `json:"call_to_actions"`
}

type action struct {
	Type    MenuType  `json:"type"`
	Title   string    `json:"title"`
	Payload string    `json:"payload,omitempty"`
	URL     string    `json:"url,omitempty"`
	WebView string    `json:"webview_height_ratio,omitempty"`
	Actions []*action `json:"call_to_actions,omitempty"`
}

// NewGetStarted creates a new GetStarted
func NewGetStarted() *GetStarted {
	return &GetStarted{
		Content: getStartedContent{Payload: "get_started"},
	}
}

// NewGreetingText creates a new GreetingText
func NewGreetingText() *GreetingText {
	return &GreetingText{
		greeting{
			Locale: "default",
			Text:   utils.GreetingTextMsg,
		},
	}
}

// NewPersistentMenu returns new persistent menu
func NewPersistentMenu() *PersistentMenu {
	return &PersistentMenu{
		Menus: []persistentMenu{
			persistentMenu{
				Locale:                "default",
				ComposerInputDisabled: false,
				Actions:               []*action{},
			},
		},
	}
}

// AddAction adds new action to menu
func (p *PersistentMenu) AddAction(a *action) {
	p.Menus[0].Actions = append(p.Menus[0].Actions, a)
}

func newWebURLAction(title, url string) *action {
	return &action{
		Type:  webURLType,
		Title: title,
		URL:   url,
	}
}

func newPostBackAction(title, payload string) *action {
	return &action{
		Type:    postBackType,
		Title:   title,
		Payload: payload,
	}
}

// GetDefaultMenu creates default persistent menu
func GetDefaultMenu() *PersistentMenu {
	pm := NewPersistentMenu()

	pm.AddAction(newPostBackAction("Latest news", "Latest news"))
	pm.AddAction(newPostBackAction("Manage Alerts", "Manage Alerts"))
	pm.AddAction(newPostBackAction("Help", "Help"))
	return pm
}

// SetupPage configure facebook page
func (mg *Messenger) SetupPage() {
	// setup get started button
	getStarted := NewGetStarted()
	_, err := mg.makeFbRequest(profilePath, "POST", getStarted)
	if err != nil {
		logger.Info("FB get started error:", err)
	}

	greetingText := NewGreetingText()
	_, err = mg.makeFbRequest(profilePath, "POST", greetingText)
	if err != nil {
		logger.Info("FB greeting setup error:", err)
	}

	menu := GetDefaultMenu()
	_, err = mg.makeFbRequest(profilePath, "POST", menu)
	if err != nil {
		logger.Info("FB menu setup error:", err)
	}
}
