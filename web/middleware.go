package web

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/utils"

	"github.com/urfave/negroni"
)

// withPostData injects post data into http.Request
func withPostData(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var data utils.Map
	defer r.Body.Close()
	// form data
	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error(err)
		}
	}
	ctx := context.WithValue(r.Context(), postDataKey, data)
	next(w, r.WithContext(ctx))
}

// auditMiddleware is a middleware to log all request in database
func auditMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Do stuff here
	tm := time.Now()
	defer func() {
		duration := time.Now().Sub(tm)
		res := w.(negroni.ResponseWriter)

		body, _ := json.Marshal(r.Context().Value(postDataKey))

		ar := models.AuditRequest{
			Path:       r.URL.Path,
			StatusCode: res.Status(),
			Proto:      r.Proto,
			Body:       string(body),
			Method:     r.Method,
			Query:      r.URL.Query().Encode(),
			UserAgent:  r.UserAgent(),
			Referrer:   r.Referer(),
			IPAddress:  r.RemoteAddr,
			Duration:   duration.String(),
			Size:       res.Size(),
		}
		// save
		ar.Save()
	}()

	// Call the next handler, which can be another middleware in the chain, or the final handler.
	next(w, r)

}
