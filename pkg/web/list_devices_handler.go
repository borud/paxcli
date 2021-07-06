package web

import (
	"fmt"
	"net/http"
)

// ListDevicesHandler lists the devices
func (s *Server) ListDevicesHandler(w http.ResponseWriter, r *http.Request) {
	devices := s.config.RingStore.ListDevices()
	w.Header().Add("Content-Type", "text/html")

	fmt.Fprintf(w, "<h1>Devices seen</h1>\n")
	fmt.Fprintf(w, "<ul>\n")
	for _, device := range devices {
		fmt.Fprintf(w, "  <li><a href=\"/devices/%s\">%s</a>", device, device)
	}
	fmt.Fprintf(w, "</ul>\n")
}
