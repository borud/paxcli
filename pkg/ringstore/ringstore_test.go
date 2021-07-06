package ringstore

import (
	"testing"
	"time"

	"github.com/borud/paxcli/pkg/model"
	"gotest.tools/assert"
)

func TestRingStore(t *testing.T) {
	ringLen := 20

	viz := New(ringLen)

	for i := 0; i < 20; i++ {
		viz.AddMeasurement(model.Measurement{
			DeviceID:             "abc",
			Timestamp:            time.Now().Add(time.Duration(i) * time.Minute),
			BluetoothDeviceCount: i,
			WIFIDeviceCount:      i,
		})
	}

	assert.Equal(t, 20, len(viz.ValuesForDeviceID("abc")))
}
