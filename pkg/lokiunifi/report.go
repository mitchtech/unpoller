package lokiunifi

import (
	"fmt"
	"strings"
	"time"

	"github.com/unpoller/unifi"
	"github.com/unpoller/unpoller/pkg/poller"
)

// LogStream contains a stream of logs (like a log file).
// This app uses one stream per log entry because each log may have different labels.
type LogStream struct {
	Labels  map[string]string `json:"stream"` // "the file name"
	Entries [][]string        `json:"values"` // "the log lines"
}

// Logs is the main logs-holding structure. This is the Loki-output format.
type Logs struct {
	Streams []LogStream `json:"streams"` // "multiple files"
}

// Report is the temporary data generated by processing events.
type Report struct {
	Start  time.Time
	Oldest time.Time
	poller.Logger
	Counts map[string]int
}

// NewReport makes a new report.
func (l *Loki) NewReport(start time.Time) *Report {
	return &Report{
		Start:  start,
		Oldest: l.last,
		Logger: l,
		Counts: make(map[string]int),
	}
}

// ProcessEventLogs loops the event Logs, matches the interface type, calls the
// appropriate method for the data, and compiles the Logs into a Loki format.
// This runs once per interval, if there was no collection error.
func (r *Report) ProcessEventLogs(events *poller.Events) *Logs {
	logs := &Logs{}

	for _, e := range events.Logs {
		switch event := e.(type) {
		case *unifi.IDS:
			r.IDS(event, logs)
		case *unifi.Event:
			r.Event(event, logs)
		case *unifi.Alarm:
			r.Alarm(event, logs)
		case *unifi.Anomaly:
			r.Anomaly(event, logs)
		default: // unlikely.
			r.LogErrorf("unknown event type: %T", e)
		}
	}

	return logs
}

func (r *Report) String() string {
	return fmt.Sprintf("%s: %d, %s: %d, %s: %d, %s: %d, Dur: %v",
		typeEvent, r.Counts[typeEvent], typeIDS, r.Counts[typeIDS],
		typeAlarm, r.Counts[typeAlarm], typeAnomaly, r.Counts[typeAnomaly],
		time.Since(r.Start).Round(time.Millisecond))
}

// CleanLabels removes any tag that is empty.
func CleanLabels(labels map[string]string) map[string]string {
	for i := range labels {
		if strings.TrimSpace(labels[i]) == "" {
			delete(labels, i)
		}
	}

	return labels
}
