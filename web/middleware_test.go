package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/urfave/negroni"

	"github.com/stretchr/testify/assert"
)

func TestWithPostData(t *testing.T) {
	assert := assert.New(t)
	// s := New()
	req := httptest.NewRequest(http.MethodPost, "/?x=y", strings.NewReader(userJSON))
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	buf := new(bytes.Buffer)
	// p := map[string]interface{}{}

	withPostData(rec, req, func(w http.ResponseWriter, r *http.Request) {
		// c := NewContext(rec, req, s)
		// p["data"] = c.PostValues()
		buf.WriteString("0")
	})

	assert.Equal("0", buf.String())
	// assert.Equal(p["data"], userJSON)
}

func TestAuditMiddleware(t *testing.T) {
	assert := assert.New(t)
	// s := New()
	req := httptest.NewRequest(http.MethodPost, "/?x=y", strings.NewReader(userJSON))
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	nw := negroni.NewResponseWriter(rec)
	buf := new(bytes.Buffer)

	auditMiddleware(nw, req, func(w http.ResponseWriter, r *http.Request) {
		buf.WriteString("0")
	})

	assert.Equal("0", buf.String())
}
