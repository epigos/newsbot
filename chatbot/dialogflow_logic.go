package chatbot

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/utils"

	dgc "github.com/mlabouardy/dialogflow-go-client"
	dgcm "github.com/mlabouardy/dialogflow-go-client/models"
)

// DialogFlowLogic logic adater that returns a response
// using https://dialogflow.com/docs
type DialogFlowLogic struct {
	bot    *Chatbot
	Client *dgc.DialogFlowClient
	state  utils.Map
}

// NewDialogFlowLogic creates a new dialogflow logic
func NewDialogFlowLogic() *DialogFlowLogic {
	err, client := dgc.NewDialogFlowClient(dgcm.Options{
		AccessToken: os.Getenv("DIALOG_FLOW_TOKEN"),
	})
	if err != nil {
		log.Fatal(err)
	}
	return &DialogFlowLogic{Client: client, state: utils.Map{}}
}

func (l *DialogFlowLogic) setChatbot(b *Chatbot) {
	l.bot = b
}

func (l *DialogFlowLogic) canProcess(s *Statement) bool {
	return true
}

// Process reads the user's input from the terminal.
func (l *DialogFlowLogic) Process(st *Statement) *Statement {
	l.bot.Logger.Debug("Using dialog flow logic")

	query := dgcm.Query{
		Query:     st.Text,
		SessionID: st.UserID,
	}
	resp, err := l.Client.QueryFindRequest(query)
	if err != nil {
		l.bot.Logger.Error(err)
	}
	l.bot.Logger.Debugf("%+v", resp)

	st.SetScore(resp.Result.Score)
	st.Meta.Set("dialog_flow_id", resp.ID)
	st.Meta.Set("timestamp", resp.Timestamp)
	st.addMessageResponseFromDialog(resp.Result.Fulfillment.Messages)

	switch action := resp.Result.Action; action {
	case utils.ActionNewsSearch:

		l.bot.Logger.Debug("Processing news search action")
		params := resp.Result.Parameters
		params["page"] = 1

		l.searchNews(st, params, 1)

	case utils.ActionNewsSearchNext:

		l.bot.Logger.Debug("Processing next news search action")
		params := l.state.Get(st.UserID, utils.Map{}).(utils.Map)
		page := params.Get("page", 1).(int)
		page = page + 1
		l.searchNews(st, params, page)

	case utils.ActionNewsSearchPrevious:

		l.bot.Logger.Debug("Processing previous news search action")
		params := l.state.Get(st.UserID, utils.Map{}).(utils.Map)
		page := params.Get("page", 1).(int)
		page = page - 1
		l.searchNews(st, params, page)

	case utils.ActionNewsSearchRepeat:

		l.bot.Logger.Debug("Processing repeat news search action")
		params := l.state.Get(st.UserID, utils.Map{}).(utils.Map)
		page := params.Get("page", 1).(int)
		l.searchNews(st, params, page)

	case utils.ActionStop:

		l.bot.Logger.Debug("Processing stop subscription action")
		if topic, ok := resp.Result.Parameters["topic"]; ok {
			l.stopSubscription(st.UserID, topic.(string))
		}

	case utils.ActionReset:

		l.bot.Logger.Debug("Processing reset subscription action")
		reply := utils.NewSubscribeMenu(st.UserID)
		st.AddResponse(reply)

	case utils.ActionTopics:

		l.bot.Logger.Debug("Processing topics list action")
		reply := utils.NewQuickReply(st.UserID, "Here are some options ⬇️")
		topics := models.GetUnsubscribedTopics(st.UserID, 5)
		for _, topic := range topics {
			reply.AddTextQuickReply(topic.Name, topic.Name)
		}
		st.AddResponse(reply)

	case utils.ActionSubscribe:

		l.bot.Logger.Debug("Processing subscribe action")
		if topic, ok := resp.Result.Parameters["topic"]; ok {
			l.subscribe(st.UserID, topic.(string))
		}

		if subs, err := models.GetUserSubscriptions(st.UserID); err == nil && len(subs) > 1 {
			break
		}
		reply := utils.NewQuickReply(st.UserID, "Do you want to subscribe to anything else?")
		reply.AddTextQuickReply("No, thanks!", "No, thanks!")
		reply.AddTextQuickReply("Other topics", "Other topics")
		st.AddResponse(reply)

	case utils.ActionManageAlerts:
		l.bot.Logger.Debug("Processing alerts action")
		if subs, err := models.GetUserSubscriptions(st.UserID); err == nil && len(subs) > 0 {
			gm := utils.NewGenericMessage(st.UserID)
			for _, sub := range subs {
				gm.AddNewElement(sub.Topic.Name, sub.Description(), "", "", sub.StopButton())
			}
			st.AddResponse(gm)
			break
		}
		st.AddTextResponse(utils.NoSubscriptionText)
		st.AddResponse(utils.NewSubscribeMenu(st.UserID))
	default:
		l.bot.Logger.Debug("Processing default action")
	}

	return st
}

func (l *DialogFlowLogic) subscribe(userID, topic string) {
	sub := models.NewSubscription(userID, topic)
	sub.Save()
}

func (l *DialogFlowLogic) stopSubscription(userID, topic string) {
	subs, _ := models.GetUserSubscriptions(userID)

	for _, sub := range subs {
		if topic == sub.Topic.Name {
			sub.Delete()
		} else if topic == "" {
			sub.Delete()
		}
	}
}

func (l *DialogFlowLogic) searchNews(st *Statement, params utils.Map, page int) error {
	// get news articles
	articles, err := models.SearchArticle(params, page)
	if err != nil {
		l.bot.Logger.Error("News search error:", err)
		return err
	}
	if len(articles) < 1 {
		return fmt.Errorf("No articles found")
	}
	// create generic message for news articles
	gm := utils.NewGenericMessage(st.UserID)
	for _, article := range articles {
		gm.AddElement(article.ToMessengerElement(st.UserID))
	}
	st.AddResponse(gm)

	// add quick replies
	cat := params.Get("category", "").(string)
	reply := utils.NewQuickReply(st.UserID, "You can view more of this news or other topics")
	reply.AddTextQuickReply("Show me more", "Show me more")
	topics := models.GetUnsubscribedTopics(st.UserID, 4)
	for _, topic := range topics {
		if topic.Name == cat {
			continue
		}
		txt := strings.Title(topic.Name) + " news"
		reply.AddTextQuickReply(txt, txt)
	}
	st.AddResponse(reply)

	// update pagination param
	params.Set("page", page)
	l.state.Set(st.UserID, params)

	return nil
}
