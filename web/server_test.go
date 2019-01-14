package web

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/epigos/newsbot/models"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"

	"github.com/joho/godotenv"
)

var srv *Server

func TestNewServer(t *testing.T) {
	assert := assert.New(t)

	assert.IsType(&Server{}, srv)
	assert.IsType(&mux.Router{}, srv.Mux)
	assert.IsType(&negroni.Negroni{}, srv.n)

	go srv.Run()
}

func TestGetRequest(t *testing.T) {
	assert := assert.New(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	srv.Mux.ServeHTTP(rec, req)

	assert.Equal(http.StatusOK, rec.Code)
	assert.Equal("Newsbot", rec.Body.String())

	req = httptest.NewRequest(http.MethodGet, "/handle", nil)
	rec = httptest.NewRecorder()
	srv.Handle("/handle", func(c *Context) *HTTPError {
		return c.WriteString("handle")
	})
	srv.Mux.ServeHTTP(rec, req)
	assert.Equal(http.StatusOK, rec.Code)
	assert.Equal("handle", rec.Body.String())
}

func TestMain(m *testing.M) {
	// load env variables
	err := godotenv.Load("../env/test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	models.Connect()

	srv = New("0.0.0.0:5059")
	m.Run()
	models.Close()
}
