package model

import "time"

// Measurement represents a pax measurement
type Measurement struct {
	DeviceID             string
	Timestamp            time.Time
	BluetoothDeviceCount int
	WIFIDeviceCount      int
	CoreTemperature      float32
	SequenceNumber       int
	SecondsUptime        int
}
