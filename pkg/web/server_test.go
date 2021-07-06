package web

import (
	"testing"

	"github.com/borud/paxcli/pkg/ringstore"
)

func TestServer(t *testing.T) {
	vizualizer := ringstore.New(100)

	server := New(Config{
		ListenAddr: ":0",
		RingStore:  vizualizer,
	})
	server.Start()
	server.Shutdown()
}
