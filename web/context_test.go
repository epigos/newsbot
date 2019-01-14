package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/epigos/newsbot/utils"

	"github.com/urfave/negroni"

	"github.com/stretchr/testify/assert"
)

var userJSON = `{"id":1,"name":"Jon Snow"}`

func TestContext(t *testing.T) {
	assert := assert.New(t)
	req := httptest.NewRequest(http.MethodGet, "/?x=y", strings.NewReader(userJSON))
	rec := httptest.NewRecorder()
	c := NewContext(rec, req, srv)

	assert.Equal(c.Server, srv)
	assert.Equal(c.Request(), req)
	assert.Equal(c.ResponseWriter, negroni.NewResponseWriter(rec))
	assert.Contains(c.Query.Get("x"), "y")
	assert.Equal(c.FormValue("x"), "y")
	var am utils.Map
	assert.Equal(c.PostValues(), &am)

	c.WriteString("test")
	assert.Equal(http.StatusOK, rec.Code)
	assert.Equal("test", rec.Body.String())
	assert.Equal(rec.Header().Get("Server"), "epigos.go")

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec = httptest.NewRecorder()
	c = NewContext(rec, req, srv)

	var a interface{}
	err := json.Unmarshal([]byte(userJSON), &a)
	assert.NoError(err)
	c.WriteJSON(a)
	assert.Equal(http.StatusOK, rec.Code)
	assert.Equal(userJSON, rec.Body.String())

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec = httptest.NewRecorder()
	c = NewContext(rec, req, srv)
	c.NotFound(errors.New("Not found"), "Not found")

	assert.Equal(http.StatusNotFound, rec.Code)
	assert.Equal("Not found", rec.Body.String())

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec = httptest.NewRecorder()
	c = NewContext(rec, req, srv)
	c.BadRequest("Validation error")
	assert.Equal(http.StatusBadRequest, rec.Code)
	assert.Equal("Validation error", rec.Body.String())

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec = httptest.NewRecorder()
	c = NewContext(rec, req, srv)
	c.NotModified()
	assert.Equal(http.StatusNotModified, rec.Code)

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec = httptest.NewRecorder()
	c = NewContext(rec, req, srv)
	c.Forbidden()
	assert.Equal(http.StatusForbidden, rec.Code)
	assert.Equal("Forbidden", rec.Body.String())

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec = httptest.NewRecorder()
	c = NewContext(rec, req, srv)
	c.Unauthorized()
	assert.Equal(http.StatusUnauthorized, rec.Code)
	assert.Equal("Unauthorized", rec.Body.String())

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec = httptest.NewRecorder()
	c = NewContext(rec, req, srv)
	c.Redirect("/?a=x")
	assert.Equal(http.StatusTemporaryRedirect, rec.Code)
}

func TestAuth(t *testing.T) {
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("test"))
	req.SetBasicAuth("user", "password")
	rec := httptest.NewRecorder()
	c := NewContext(rec, req, srv)

	user, pwd, err := c.GetBasicAuth()
	if assert.NoError(err) {
		assert.Equal(user, "user")
		assert.Equal(pwd, "password")
	}
}
