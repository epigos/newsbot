package web

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/epigos/newsbot/utils"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

const (
	postDataKey key = "Data"
)

type key string

// Context request context
type Context struct {
	negroni.ResponseWriter
	r      *http.Request
	Server *Server
	Params map[string]string
	Query  url.Values
}

// NewContext creates new instance of context
func NewContext(w http.ResponseWriter, req *http.Request, s *Server) *Context {
	p := mux.Vars(req)
	sw := negroni.NewResponseWriter(w)
	ctx := &Context{sw, req, s, p, req.URL.Query()}
	return ctx
}

// Request returns the http request
func (ctx *Context) Request() *http.Request {
	return ctx.r
}

// WriteString writes string data into the response object.
func (ctx *Context) WriteString(content string) *HTTPError {
	ctx.setDefaultHeaders()
	ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	// set the default content-type
	ctx.WriteHeader(http.StatusOK)

	if _, err := ctx.ResponseWriter.Write([]byte(content)); err != nil {
		return serverError(err)
	}
	return nil
}

// WriteJSON writes string data into the response object.
func (ctx *Context) WriteJSON(data interface{}) *HTTPError {
	ctx.setDefaultHeaders()
	json, err := json.Marshal(data)
	if err != nil {
		return serverError(err)
	}
	ctx.ContentType("application/json", true)
	// set the default content-type
	ctx.WriteHeader(http.StatusOK)
	if _, err := ctx.ResponseWriter.Write(json); err != nil {
		return serverError(err)
	}
	return nil
}

func (ctx *Context) setDefaultHeaders() {
	ctx.ContentType("text/html; charset=utf-8", true)
	ctx.SetHeader("Server", "epigos.go", true)
	ctx.SetHeader("Date", utils.WebTime(time.Now().UTC()), true)
	ctx.SetHeader("Access-Control-Allow-Origin", "*", true)
	ctx.SetHeader("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding", true)
	ctx.SetHeader("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS", true)
}

// Redirect is a helper method for 3xx redirects.
func (ctx *Context) Redirect(url string) {
	ctx.ResponseWriter.Header().Set("Location", url)
	ctx.ResponseWriter.WriteHeader(http.StatusTemporaryRedirect)
	ctx.ResponseWriter.Write([]byte("Redirecting to: " + url))
}

//ServerError writes a 500 HTTP response
func (ctx *Context) ServerError(err error) *HTTPError {
	return serverError(err)
}

//BadRequest writes a 400 HTTP response
func (ctx *Context) BadRequest(s string) *HTTPError {
	return badRequestError(s)
}

// NotModified writes a 304 HTTP response
func (ctx *Context) NotModified() {
	ctx.WriteError(http.StatusNotModified, "")
}

//Unauthorized writes a 401 HTTP response
func (ctx *Context) Unauthorized() {
	ctx.WriteError(http.StatusUnauthorized, "Unauthorized")
}

//Forbidden writes a 403 HTTP response
func (ctx *Context) Forbidden() {
	ctx.WriteError(http.StatusForbidden, "Forbidden")
}

// NotFound writes a 404 HTTP response
func (ctx *Context) NotFound(err error, message string) *HTTPError {
	return notFoundError(err, message)
}

// ContentType sets the Content-Type header for an HTTP response.
// For example, ctx.ContentType("json") sets the content-type to "application/json"
// If the supplied value contains a slash (/) it is set as the Content-Type
// verbatim. The return value is the content type as it was
// set, or an empty string if none was found.
func (ctx *Context) ContentType(val string, unique bool) {
	var ctype string
	if strings.ContainsRune(val, '/') {
		ctype = val
	} else {
		if !strings.HasPrefix(val, ".") {
			val = "." + val
		}
		ctype = mime.TypeByExtension(val)
	}
	if ctype != "" {
		ctx.SetHeader("Content-Type", ctype, unique)
	}
}

// SetHeader sets a response header. If `unique` is true, the current value
// of that header will be overwritten . If false, it will be appended.
func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.Header().Set(hdr, val)
	} else {
		ctx.Header().Add(hdr, val)
	}
}

// SetCookie adds a cookie header to the response.
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	ctx.SetHeader("Set-Cookie", cookie.String(), false)
}

// FormValue returns the form field value for the provided name.
func (ctx *Context) FormValue(name string) string {
	return ctx.Request().FormValue(name)
}

// WriteError writes error response
func (ctx *Context) WriteError(status int, m string) {
	ctx.setDefaultHeaders()
	ctx.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte(m))
}

// PostValues returns post data in request
func (ctx *Context) PostValues() *utils.Map {
	a := ctx.Request().Context().Value(postDataKey)
	if a == nil {
		var a utils.Map
		defer ctx.Request().Body.Close()
		if ctx.Request().Method == http.MethodPost {
			err := json.NewDecoder(ctx.Request().Body).Decode(&a)
			if err != nil {
				logger.Error(err)
			}
		}
		return &a
	}
	p := a.(utils.Map)
	return &p
}

// GetBasicAuth returns the decoded user and password from the context's
// 'Authorization' header.
func (ctx *Context) GetBasicAuth() (string, string, error) {
	if len(ctx.Request().Header["Authorization"]) == 0 {
		return "", "", errors.New("No Authorization header provided")
	}
	authHeader := ctx.Request().Header["Authorization"][0]
	authString := strings.Split(string(authHeader), " ")
	if authString[0] != "Basic" {
		return "", "", errors.New("Not Basic Authentication")
	}
	decodedAuth, err := base64.StdEncoding.DecodeString(authString[1])
	if err != nil {
		return "", "", err
	}
	authSlice := strings.Split(string(decodedAuth), ":")
	if len(authSlice) != 2 {
		return "", "", errors.New("Error delimiting authString into username/password. Malformed input: " + authString[1])
	}
	return authSlice[0], authSlice[1], nil
}

// GetParams returns url params in request
func (ctx *Context) GetParams() map[string]string {
	return ctx.Params
}

// GetParam returns url param by key in request
func (ctx *Context) GetParam(key string) string {
	return ctx.Params[key]
}

// GetQuery returns url query in request
func (ctx *Context) GetQuery() url.Values {
	return ctx.Query
}
