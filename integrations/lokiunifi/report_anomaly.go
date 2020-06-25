package lokiunifi

import (
	"strconv"

	"github.com/unifi-poller/unifi"
)

const typeAnomaly = "anomaly"

// Anomaly stores a structured Anomaly for batch sending to Loki.
func (r *Report) Anomaly(event *unifi.Anomaly) {
	if event.Datetime.Before(*r.Last) {
		return
	}

	r.Counts[typeAnomaly]++ // increase counter and append new log line.
	r.Streams = append(r.Streams, LogStream{
		Entries: [][]string{{strconv.FormatInt(event.Datetime.UnixNano(), 10), event.Anomaly}},
		Labels: CleanLabels(map[string]string{
			"application": "unifi_anomaly",
			"source":      event.SourceName,
			"site_name":   event.SiteName,
			"device_mac":  event.DeviceMAC,
		}),
	})
}
