package messenger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/epigos/newsbot/web"

	"github.com/stretchr/testify/assert"
)

var (
	entry = `{"entry": [{"id": "1251562161607", "time": 1527904873637, "messaging": [{"recipient": {"id": "1251562161607"}, "sender": {"id": "1403078893046"}, "timestamp": 1527904863157, "message": {"mid": "mid.$cAARySlSHk35p739LtVjvjoDMaVaX", "seq": 18934, "text": "hi"}}]}], "object": "page"}`
)

func TestFacebookRequest(t *testing.T) {
	assert := assert.New(t)
	req := httptest.NewRequest("POST", "/facebook", strings.NewReader(entry))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	d := map[string]*FacebookRequest{}

	srv.Handle("/facebook", func(c *web.Context) *web.HTTPError {
		fb, err := mg.DecodeRequest(c)
		if err != nil {
			c.WriteError(http.StatusInternalServerError, err.Error())
		}
		d["fb"] = fb
		return c.WriteString("Message received")
	})
	srv.Mux.ServeHTTP(rec, req)
	assert.Equal(http.StatusOK, rec.Code)
	assert.Equal("Message received", rec.Body.String())
	fb := d["fb"]
	assert.Equal(fb.Entry[0].ID, "1251562161607")
	assert.Equal(fb.Entry[0].Messaging[0].String(), "From: 1403078893046, Text: hi")
}

func TestFacebookResponse(t *testing.T) {
	assert := assert.New(t)
	// fs will mock up fb messenger server
	fs := getFbServer()
	defer fs.Close()
	resp, err := http.Get(fs.URL)
	assert.NoError(err)
	fbRes, err := mg.decodeResponse(resp)
	assert.NoError(err)
	assert.Equal(fbRes.MessageID, mid)
	assert.Equal(fbRes.RecipientID, rid)
	s := fmt.Sprintf("FB response: %s %v", mid, rid)
	assert.Equal(fbRes.String(), s)

	fs = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := rawFBResponse{
			RecipientID: rid,
			MessageID:   mid,
			Error:       &FacebookError{400, "1", "Invalid message format", "error"},
		}
		b, _ := json.Marshal(rec)
		w.Write(b)
	}))
	defer fs.Close()
	resp, err = http.Get(fs.URL)
	assert.NoError(err)
	fbRes, err = mg.decodeResponse(resp)
	assert.NotEqual(fbRes.MessageID, mid)
	assert.Error(err)
	assert.Equal(err, fmt.Errorf("FB Error: Type error: Invalid message format; FB trace ID: 1"))
}
