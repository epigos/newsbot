package web

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/epigos/newsbot/models"
	"github.com/epigos/newsbot/utils"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
)

var (
	logger = utils.NewLogger("web")
)

// Server main web application
type Server struct {
	Config utils.Map
	Mux    *mux.Router
	n      *negroni.Negroni
	Logger *utils.Logger
}

// New creates a new server
func New(h string) *Server {

	cfg := utils.Map{"host": h}
	app := &Server{
		Config: cfg,
		Logger: logger,
	}
	// add middlewares
	mux := mux.NewRouter()
	n := negroni.Classic() // Includes some default middlewares
	n.Use(negroni.HandlerFunc(withPostData))
	n.Use(negroni.HandlerFunc(auditMiddleware))
	n.UseHandler(mux)
	// add router
	app.Mux = mux
	app.n = n
	// add default handlers
	app.Mux.HandleFunc(`/*`, http.NotFound)
	app.ConfigureRoute()
	return app
}

//Run new web application
func (s *Server) Run() {
	host := s.Config.Get("host", nil).(string)

	certcache := NewDatastoreCertCache(models.DS.Client)

	certManager := autocert.Manager{
		Cache:      certcache,
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(os.Getenv("HOST_NAME")),
		Email:      os.Getenv("ADMIN_EMAIL"),
	}

	tlsConfig := &tls.Config{
		Rand:           rand.Reader,
		Time:           time.Now,
		NextProtos:     []string{http2.NextProtoTLS, "http/1.1"},
		MinVersion:     tls.VersionTLS12,
		GetCertificate: certManager.GetCertificate,
	}

	srv := &http.Server{
		Addr: host,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      s.n, // Pass our instance of gorilla/mux in.
	}

	if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == "prod" {
		srv.Addr = ":443"
		srv.TLSConfig = tlsConfig
		s.Logger.Infof("Listening on https://%v/", host)
		go http.ListenAndServe(host, certManager.HTTPHandler(nil))
		go srv.ListenAndServeTLS("", "")
	} else {
		// Run our server in a goroutine so that it doesn't block.
		go func() {
			s.Logger.Infof("Listening on http://%v/", host)
			if err := srv.ListenAndServe(); err != nil {
				s.Logger.Critical(err)
			}
		}()
	}

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	s.Logger.Info("shutting down")
	os.Exit(0)
}

// Handle custom routes
func (s *Server) Handle(path string, handler httpHandler, methods ...string) {
	s.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, s)

		if err := handler(ctx); err != nil {
			http.Error(w, err.Message, err.Code)
		}

	}).Methods(methods...).Host(os.Getenv("HOST_NAME"))
}

// Get handles GET request to url path
func (s *Server) Get(path string, handler httpHandler) {
	s.Handle(path, handler, http.MethodGet)
}

// Post handles POST request to url path
func (s *Server) Post(path string, handler httpHandler) {
	s.Handle(path, handler, http.MethodPost)
}

// Put handles PUT request to url path
func (s *Server) Put(path string, handler httpHandler) {
	s.Handle(path, handler, http.MethodPut)
}
