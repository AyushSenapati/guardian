package server

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AyushSenapati/guardian/config"
	"github.com/AyushSenapati/guardian/lib/middleware"
	"github.com/AyushSenapati/guardian/lib/proxy"
	"github.com/AyushSenapati/guardian/lib/router"
	"github.com/AyushSenapati/guardian/lib/service"
)

// Server defines the core DS of Guardian server
type Server struct {
	server        *http.Server
	globalConfig  *config.Specification
	Register      *proxy.Register
	ServiceLoader *service.Loader
	stopChan      chan struct{}
}

// NewServerWithConfig takes the config specification and returns a server obj
func NewServerWithConfig(config *config.Specification) *Server {
	return &Server{
		globalConfig: config,
		stopChan:     make(chan struct{}, 1),
	}
}

// Start will start the gateway service
func (s *Server) Start(ctx context.Context, svcDefinitionFname string) {

	// this triggers the server close once
	// the parent context gets canceled
	go func() {
		<-ctx.Done()
		log.Println(
			"Guardian >> Oopps! Gotta go. till then manage your services",
			html.UnescapeString("&#"+strconv.Itoa(128517)+";"),
		)
		s.Close()
	}()

	// Create a router interface
	r := s.CreateRouter()
	// Register the router interface in the proxy register
	s.Register = proxy.NewRegister(r)
	// Register the proxy register in service loader
	s.ServiceLoader = service.NewLoader(s.Register)

	go func() {
		if err := s.startHTTPServer(r); err == http.ErrServerClosed {
			log.Println(
				"Guardian >> mmm, lemme wait till your active connections're closed...")
		} else {
			log.Fatal(err)
		}
	}()

	// routes can be loaded even after the server has started.
	// But till the routes are loaded it is obvious
	// accessing those routes will give 404 NotFound
	// r.LoadRouterDefinitions("definitions.json")
	definitions := s.ServiceLoader.LoadServiceDefinitions(svcDefinitionFname)
	s.ServiceLoader.RegisterServices(definitions)
}

// Wait will wait till any signal is
// received on server stop channel to stop the server
func (s *Server) Wait(ctx context.Context) {
	log.Printf("info: PPID:%d, PID:%d", os.Getppid(), os.Getpid())
	<-s.stopChan
	time.Sleep(2 * time.Second)
}

// Close shutdowns the server gracefully
func (s *Server) Close() error {
	defer close(s.stopChan)

	ctx, cancel := context.WithTimeout(
		context.Background(), time.Duration(s.globalConfig.GraceTimeout)*time.Second)
	defer cancel()

	go func(ctx context.Context) {
		<-ctx.Done()
		if ctx.Err() == context.Canceled {
			return
		} else if ctx.Err() == context.DeadlineExceeded {
			log.Println(
				"Guardian >> F**k! lost my patience. force terminating connections",
			)
		}
	}(ctx)

	return s.server.Shutdown(ctx)
}

// it creates the http.Server instance, listens and serves http requests
func (s *Server) startHTTPServer(r router.Router) error {
	addr := fmt.Sprintf(":%v", s.globalConfig.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(s.globalConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.globalConfig.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.globalConfig.IdleTimeout) * time.Second,
	}

	log.Println("srv: will listen on", addr)

	return s.server.ListenAndServe()
}

// CreateRouter returns a router interface
func (s *Server) CreateRouter() router.Router {
	r := router.NewHTTPServeMux()

	// RequestID must be the first middleware to be registered
	// if it is configured, so that all other handlers can access requestID
	// from the context to log the errors in case any
	if s.globalConfig.AddReqID {
		r.Use(middleware.RequestID)
	}

	return r
}
