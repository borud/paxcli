package main

import (
	"fmt"
	"log"
	"os"

	"github.com/borud/paxcli/pkg/spanlisten"
	"github.com/jessevdk/go-flags"
)

var opt struct {
	Token        string `long:"token" env:"SPAN_API_TOKEN" description:"Span API Token" required:"yes"`
	CollectionID string `long:"collection" description:"Span Collection ID for PAX counters" default:"17dlb1hl0l800a"`
}

var parser = flags.NewParser(&opt, flags.Default)

func init() {
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}

func main() {
	spanListen := spanlisten.New(opt.Token, opt.CollectionID)
	err := spanListen.Start()
	if err != nil {
		log.Fatalf("error starting spanlistener: %v", err)
	}

	for m := range spanListen.Measurements() {
		fmt.Printf("%s device='%s' bluetooth=%d wifi=%d\n", m.Timestamp.Format("2006-01-02 15:04:05"), m.DeviceID, m.BluetoothDeviceCount, m.WIFIDeviceCount)
	}
}
