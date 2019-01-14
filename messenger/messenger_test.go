package messenger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"github.com/epigos/newsbot/chatbot"
	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/web"
	"strings"
	"testing"
	"time"

	"github.com/icrowley/fake"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	rid  = "12123213123"
	mid  = "mid00000TEST00000TEST00000TEST"
	user *models.User
	srv  *web.Server
	mg   *Messenger
)

func getChatServer() (*httptest.Server, *Messenger) {
	fs := getFbServer()
	defer fs.Close()

	ch := chatbot.New("Test")
	msng := New(ch)
	// ts is our test chatbot
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := web.NewContext(w, r, srv)
		msng.ServeHTTP(c)
	}))

	return ts, msng
}

func getFbServer() *httptest.Server {
	// fs will mock up fb messenger server
	fs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := FacebookResponse{
			RecipientID: rid,
			MessageID:   mid,
		}
		b, _ := json.Marshal(rec)
		w.Write(b)
	}))
	TestURL = fs.URL
	return fs
}

func TestVerify(t *testing.T) {
	assert := assert.New(t)
	ts, mg := getChatServer()
	defer ts.Close()
	challenge := "1122334455"
	verifyReq := ts.URL + "/?test=1&hub.mode=subscribe&hub.challenge=" + challenge + "&hub.verify_token=" + mg.VerifyToken
	resp, err := http.Get(verifyReq)
	assert.NoError(err)
	defer resp.Body.Close()
	s, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(string(s), challenge)
}

func TestBuildUrl(t *testing.T) {
	assert := assert.New(t)
	ts, mg := getChatServer()
	defer ts.Close()

	url := mg.buildURL(messagesPath)
	u := fmt.Sprintf("%s/me/%s?access_token=%s", TestURL, messagesPath, mg.AccessToken)
	assert.Equal(url, u)

	url = mg.buildURL("messagesPaths")
	un := fmt.Sprintf("%s/messagesPaths?access_token=%s", TestURL, mg.AccessToken)
	assert.Equal(url, un)
	assert.NotEqual(url, u)

	ua := fmt.Sprintf("%s/me/%s?access_token=%s", apiURL, messagesPath, mg.AccessToken)
	TestURL = ""
	url = mg.buildURL(messagesPath)
	assert.Equal(url, ua)
	assert.NotEqual(url, u)
}

func TestFbRequest(t *testing.T) {
	assert := assert.New(t)
	fs := getFbServer()
	defer fs.Close()
	resp, err := mg.makeFbRequest(messagesPath, "POST", strings.NewReader(entry))
	assert.NoError(err)
	fbRes, err := mg.decodeResponse(resp)
	assert.NoError(err)
	assert.Equal(fbRes.MessageID, mid)
	assert.Equal(fbRes.RecipientID, rid)
}

func TestHandleMessage(t *testing.T) {
	assert := assert.New(t)
	ts, mg := getChatServer()
	defer ts.Close()
	mg.Handler = &messageTestHandler{mg.Handler}
	req := httptest.NewRequest("POST", "/facebook", strings.NewReader(entry))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.Handle("/facebook", mg.ServeHTTP)
	srv.Mux.ServeHTTP(rec, req)
	go mg.Listen()
	time.Sleep(500 * time.Millisecond)
	assert.Equal(http.StatusOK, rec.Code)
	assert.Equal("Message received", rec.Body.String())
}

type messageTestHandler struct {
	MessageHandler
}

func (h *messageTestHandler) ProcessMessage(m *messaging) {
	log.Println(m)
}

// func (h *messageTestHandler) ProcessDelivery(m *messaging) {

// }
// func (h *messageTestHandler) ProcessPostback(m *messaging) {

// }

func TestGetSenderProfile(t *testing.T) {
	assert := assert.New(t)
	id := fake.Characters()
	nfs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := map[string]interface{}{
			"id":          id,
			"first_name":  fake.FirstName(),
			"last_name":   fake.LastName(),
			"profile_pic": fake.DomainName(),
			"locale":      "en_US",
			"timezone":    0,
			"gender":      fake.Gender(),
		}
		b, _ := json.Marshal(rec)
		w.Write(b)
	}))
	TestURL = nfs.URL
	bu := mg.GetSenderProfile(id)
	assert.Equal(bu.ID, id)
}

func TestMain(m *testing.M) {
	// load env variables
	err := godotenv.Load("../env/test.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ch := chatbot.New("Test")
	mg = New(ch)

	srv = web.New("0.0.0.0:5050")
	models.Connect()

	user = models.NewUser(fake.Characters(), fake.FirstName(), fake.LastName(), fake.DomainName(), fake.Language(), fake.Gender(), 0)

	m.Run()
	models.Close()
}
