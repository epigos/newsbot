package chatbot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/utils"
)

// PostBackLogic logic to processs postback
type PostBackLogic struct {
	bot     *Chatbot
	Actions []string
	regex   *regexp.Regexp
}

// NewPostBackLogic returns new PostBackLogic
func NewPostBackLogic() *PostBackLogic {
	actions := []string{
		utils.PostBackGetStarted,
		utils.PostBackGetSummary,
	}
	regex := regexp.MustCompile(fmt.Sprintf(`%s`, strings.Join(actions, "|")))
	return &PostBackLogic{Actions: actions, regex: regex}
}

func (l *PostBackLogic) setChatbot(b *Chatbot) {
	l.bot = b
}

func (l *PostBackLogic) canProcess(s *Statement) bool {
	return l.regex.MatchString(s.Text)
}

// Process reads the user's input from the terminal.
func (l *PostBackLogic) Process(st *Statement) *Statement {
	l.bot.Logger.Debug("Using postback logic")

	switch st.Text {
	case utils.PostBackGetStarted:
		l.bot.Logger.Debugf("Processing %s", utils.PostBackGetStarted)

		st.AddTextResponse(utils.GetStartedMsg)
		st.AddTextResponse(utils.SubscribeText)

		reply := utils.NewSubscribeMenu(st.UserID)
		st.AddResponse(reply)

	case utils.PostBackGetSummary:
		l.bot.Logger.Debugf("Processing %s", utils.PostBackGetSummary)

		article, err := models.GetArticle(st.Payload)
		if err != nil {
			l.bot.Logger.Error("Article summary:", err)
		} else {
			for _, sumr := range article.Summary {
				st.AddTextResponse(sumr)
			}
			// save user action and update article score
			go func(s *Statement, a *models.Article) {
				ua := models.NewUserAction(s.UserID, a.Key(), "summary")
				ua.Save()
				a.Score += utils.SummaryScore
				a.Save()
			}(st, article)
		}
	default:
		l.bot.Logger.Debugf("Default post back: %+v", st.Text)
	}

	return st
}
