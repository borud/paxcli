package ringstore

import (
	"container/ring"
	"log"

	"github.com/borud/paxcli/pkg/model"
)

// RingStore is a ringbuffer storage for measurements
type RingStore struct {
	ringLen int
	devices map[string]*ring.Ring
}

// New creates a new RingStore
func New(ringLen int) *RingStore {
	return &RingStore{
		ringLen: ringLen,
		devices: make(map[string]*ring.Ring),
	}
}

// AddMeasurement adds a model.Measurement to the RingStore
func (v *RingStore) AddMeasurement(m model.Measurement) {
	r, ok := v.devices[m.DeviceID]
	if !ok {
		r = ring.New(v.ringLen)
	}

	r.Value = m
	v.devices[m.DeviceID] = r.Next()
}

// ListDevices lists the devices we have seen
func (v *RingStore) ListDevices() []string {
	devices := []string{}
	for k := range v.devices {
		devices = append(devices, k)
	}
	return devices
}

// ValuesForDeviceID returns an array of mode.Measurement in ascending order of time added
func (v *RingStore) ValuesForDeviceID(deviceID string) []model.Measurement {
	r, ok := v.devices[deviceID]
	if !ok {
		return nil
	}

	var values []model.Measurement

	r.Do(func(p interface{}) {
		if p == nil {
			return
		}
		measurement, ok := p.(model.Measurement)
		if !ok {
			log.Printf("error, not measurement")
		}
		values = append(values, measurement)
	})

	return values
}
