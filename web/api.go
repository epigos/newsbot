package web

import (
	"strconv"

	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/utils"
)

// facebookUsersView handler to retrieve facebook users
func facebookUsersView(ctx *Context) *HTTPError {
	// mu := models.BotUser{}
	// users, err := mu.FindAll()
	// if err != nil {
	// 	ctx.WriteError(http.StatusInternalServerError, err.Error())
	// }
	// a := map[string][]*models.BotUser{"bot_users": users}
	// ctx.WriteJSON(&a)
	return nil
}

func searchAPI(ctx *Context) *HTTPError {
	q := ctx.GetQuery()

	m := utils.Map{
		"category": q.Get("topic"),
		"keyword":  q.Get("q"),
		"source":   q.Get("src"),
	}
	p, _ := strconv.Atoi(q.Get("page"))

	var articles []*models.Article
	var err error

	if articles, err = models.SearchArticle(m, p); err != nil {
		return ctx.ServerError(err)
	}
	return ctx.WriteJSON(articles)
}
