package web

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	chart "github.com/wcharczuk/go-chart/v2"
)

// GraphDeviceHandler graphs one device
func (s *Server) GraphDeviceHandler(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := mux.Vars(r)["deviceID"]
	if !ok {
		http.Error(w, "deviceID not specified", http.StatusBadRequest)
		return
	}

	values := s.config.RingStore.ValuesForDeviceID(deviceID)
	if values == nil {
		http.Error(w, "deviceID not found", http.StatusNotFound)
		return
	}

	times := []time.Time{}
	wifiCounts := []float64{}
	btCounts := []float64{}
	coreTemps := []float64{}

	for _, m := range values {
		times = append(times, m.Timestamp)
		wifiCounts = append(wifiCounts, float64(m.WIFIDeviceCount))
		btCounts = append(btCounts, float64(m.BluetoothDeviceCount))
		coreTemps = append(coreTemps, float64(m.CoreTemperature))
	}

	graph := chart.Chart{
		Title: "Counters",
		TitleStyle: chart.Style{
			Padding: chart.Box{
				Top:    10,
				Left:   10,
				Right:  10,
				Bottom: 10,
			},
		},
		Width:  1200,
		Height: 800,
		Background: chart.Style{
			Hidden: false,
			Padding: chart.Box{
				Top:    10,
				Left:   10,
				Right:  10,
				Bottom: 10,
			},
		},

		Series: []chart.Series{
			// Bluetooth series
			chart.TimeSeries{
				Name:    "bluetooth",
				XValues: times,
				YValues: btCounts,
			},
			// Wifi series
			chart.TimeSeries{
				Name:    "wifi",
				XValues: times,
				YValues: wifiCounts,
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	w.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, w)
}
