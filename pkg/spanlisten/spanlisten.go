package spanlisten

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/borud/paxcli/pkg/apipb"
	"github.com/borud/paxcli/pkg/model"
	"github.com/lab5e/go-spanapi/v4"
	"github.com/lab5e/go-spanapi/v4/apitools"
	"google.golang.org/protobuf/proto"
)

// SpanListener listens to a given collection on Span
type SpanListener struct {
	Token            string
	CollectionID     string
	measurementCh    chan model.Measurement
	cancel           context.CancelFunc
	ctx              context.Context
	shutdownComplete sync.WaitGroup
}

// New creates a new SpanListener instance
func New(token string, collectionID string) *SpanListener {
	return &SpanListener{
		Token:         token,
		CollectionID:  collectionID,
		measurementCh: make(chan model.Measurement),
	}
}

// Start fires up the Spanlistener
func (s *SpanListener) Start() error {
	config := spanapi.NewConfiguration()
	config.Debug = true

	s.ctx, s.cancel = context.WithCancel(apitools.ContextWithAuth(s.Token))
	ds, err := apitools.NewCollectionDataStream(s.ctx, config, s.CollectionID)
	if err != nil {
		return fmt.Errorf("unable to open CollectionDataStream: %v", err)
	}

	// Start goroutine running readDataStream() function
	go s.readDataStream(ds)

	return nil
}

// Stop listener
func (s *SpanListener) Stop() {
	if s.cancel != nil {
		s.cancel()
		s.shutdownComplete.Wait()
	}
}

// Measurements returns a channel that streams apipb.PAXMessage
func (s *SpanListener) Measurements() <-chan model.Measurement {
	return s.measurementCh
}

func (s *SpanListener) readDataStream(ds apitools.DataStream) {
	defer ds.Close()

	// Signal that we have started
	s.shutdownComplete.Add(1)

	log.Printf("connected to Span")
	for {
		msg, err := ds.Recv()
		if err != nil {
			log.Fatalf("error reading message: %v", err)
		}

		// We only care about messages containing data
		if *msg.Type != "data" {
			continue
		}

		// base64 decode the payload to a string
		bytePayload, err := base64.StdEncoding.DecodeString(*msg.Payload)
		if err != nil {
			log.Fatalf("unable to decode payload: %v", err)
		}

		// decode bytePayload as protobuffer
		var pb apipb.PAXMessage
		err = proto.Unmarshal(bytePayload, &pb)
		if err != nil {
			log.Fatalf("unable to unmarshal protobuffer: %v", err)
		}

		timeMS, err := strconv.ParseInt(*msg.Received, 10, 64)
		if err != nil {
			log.Printf("error parsing '%s' as timestamp: %v", *msg.Received, err)
			continue
		}
		timeStamp := time.Unix(0, timeMS*int64(time.Millisecond))

		s.measurementCh <- model.Measurement{
			DeviceID:             *msg.Device.DeviceId,
			Timestamp:            timeStamp,
			BluetoothDeviceCount: int(pb.BluetoothDeviceCount),
			WIFIDeviceCount:      int(pb.WifiDeviceCount),
			CoreTemperature:	  float32(pb.CoreTemperature),
			SequenceNumber: 	  int(pb.SequenceNumber),
			SecondsUptime:		  int(pb.SecondsUptime),
		}

		if s.ctx.Err() == context.Canceled {
			log.Printf("shutting down spanlistener")
			close(s.measurementCh)
			s.shutdownComplete.Done()
			return
		}
	}
}
