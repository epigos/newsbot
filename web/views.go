package web

import (
	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/utils"
)

// HomeView handler for home page
func homeView(ctx *Context) *HTTPError {
	return ctx.WriteString("Newsbot")
}

// newsRedirectView handler for news
func articleRedirectView(ctx *Context) *HTTPError {

	articleID := ctx.GetParam("articleID")
	userID := ctx.GetParam("userID")

	akey := models.DS.DecodeKey(articleID)

	go func(uid, aid string) {
		if article, err := models.GetArticle(aid); err == nil {
			ua := models.NewUserAction(uid, article.Key(), "view")
			ua.Save()
			article.Score += utils.ViewScore
			article.Save()
		}
	}(userID, akey.Name)

	ctx.Redirect(akey.Name)
	return nil
}

func articleView(ctx *Context) *HTTPError {
	articleID := ctx.GetParam("articleID")

	var article *models.Article
	var err error

	if article, err = models.GetArticle(articleID); err != nil {
		return ctx.NotFound(err, "Article does not exist")
	}
	return ctx.WriteJSON(article)
}
