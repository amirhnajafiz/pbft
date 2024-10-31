package metrics

import "time"

// Metrics is a module that stores a node performance.
type Metrics struct {
	latency    float64 // time spent for processing a single transaction in microsecond
	throughput float64 // transactions per second parameter
	records    int     // hold the number of records
}

// NewMetrics returns a new metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{
		latency:    float64(0),
		throughput: float64(0),
		records:    0,
	}
}

// Observe get's a duration time and updates system metrics.
func (m *Metrics) Observe(duration time.Duration) {
	// get microseconds
	value := duration.Nanoseconds()

	// calculate values
	var throughput float64
	if value != 0 {
		throughput = 1000000000 * float64(1/float64(value))
	} else {
		throughput = 0
	}

	// update records
	m.records++

	// calculate the average value
	m.latency += float64(value / 1000)
	m.throughput += throughput
}

// GetValues is used to export the current metrics. (latency, throughput)
func (m *Metrics) GetValues() (float64, float64) {
	if m.records == 0 {
		return m.latency, m.throughput
	}

	return m.latency / float64(m.records), m.throughput / float64(m.records)
}
