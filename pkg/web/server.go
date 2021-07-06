package web

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/borud/paxcli/pkg/ringstore"
	"github.com/gorilla/mux"
)

// Server is the webserver
type Server struct {
	config        Config
	serverStarted sync.WaitGroup
	serverStopped sync.WaitGroup
	cancel        context.CancelFunc
	ctx           context.Context
}

// Config for webserver
type Config struct {
	ListenAddr string
	RingStore  *ringstore.RingStore
}

// New webserver
func New(config Config) *Server {
	return &Server{
		config:        config,
		serverStarted: sync.WaitGroup{},
		serverStopped: sync.WaitGroup{},
	}
}

// Start webserver
func (s *Server) Start() error {
	mux := mux.NewRouter()

	mux.HandleFunc("/", s.ListDevicesHandler).Methods("GET")
	mux.HandleFunc("/devices/{deviceID}", s.GraphDeviceHandler).Methods("GET")

	s.ctx, s.cancel = context.WithCancel(context.Background())

	httpServer := &http.Server{
		Handler: mux,
	}

	s.serverStarted.Add(1)
	s.serverStopped.Add(1)

	listener, err := net.Listen("tcp", s.config.ListenAddr)
	if err != nil {
		return err
	}

	go func() {
		<-s.ctx.Done()
		httpServer.Shutdown(context.Background())
		s.serverStopped.Done()
		log.Printf("shut down webserver")
	}()

	go func() {
		s.serverStarted.Done()
		log.Printf("webserver listening to %s", listener.Addr().String())
		err := httpServer.Serve(listener)
		if err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	s.serverStarted.Wait()
	return nil
}

// Shutdown webserver
func (s *Server) Shutdown() {
	if s.cancel != nil {
		s.cancel()
		s.serverStopped.Wait()
	}
}
