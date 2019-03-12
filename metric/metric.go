package metric

import (
	"net/http"
	"os"

	"github.com/juju/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/projecteru2/yavirt/util"
)

var (
	DefaultLabels = []string{"host"}

	MetricHeartbeatCount = "yavirt_heartbeat_total"
	MetricErrorCount     = "yavirt_error_total"

	metr *Metrics
)

func init() {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	metr = New(host)
	metr.RegisterCounter(MetricErrorCount, "yavirt errors", nil)
	metr.RegisterCounter(MetricHeartbeatCount, "yavirt heartbeats", nil)
}

type Metrics struct {
	host     string
	counters map[string]*prometheus.CounterVec
	gauges   map[string]*prometheus.GaugeVec
}

func New(host string) *Metrics {
	return &Metrics{
		host:     host,
		counters: map[string]*prometheus.CounterVec{},
		gauges:   map[string]*prometheus.GaugeVec{},
	}
}

func (m *Metrics) RegisterCounter(name, desc string, labels []string) error {
	var col = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: desc,
		},
		util.MergeStrings(labels, DefaultLabels),
	)

	if err := prometheus.Register(col); err != nil {
		return errors.Trace(err)
	}

	m.counters[name] = col

	return nil
}

func (m *Metrics) RegisterGauge(name, desc string, labels []string) error {
	var col = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: desc,
		},
		util.MergeStrings(labels, DefaultLabels),
	)

	if err := prometheus.Register(col); err != nil {
		return errors.Trace(err)
	}

	m.gauges[name] = col

	return nil
}

func (m *Metrics) Incr(name string, labels map[string]string) error {
	var col, exists = m.counters[name]
	if !exists {
		return errors.Errorf("collector %s not found", name)
	}

	labels = m.appendLabel(labels, "host", m.host)

	col.With(labels).Inc()

	return nil
}

func (m *Metrics) Store(name string, value float64, labels map[string]string) error {
	var col, exists = m.gauges[name]
	if !exists {
		return errors.Errorf("collector %s not found", name)
	}

	labels = m.appendLabel(labels, "host", m.host)

	col.With(labels).Set(value)

	return nil
}

func (m *Metrics) appendLabel(labels map[string]string, key, value string) map[string]string {
	if labels != nil {
		labels["host"] = m.host
	} else {
		labels = map[string]string{"host": m.host}
	}
	return labels
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func IncrError() {
	Incr(MetricErrorCount, nil)
}

func IncrHeartbeat() {
	Incr(MetricHeartbeatCount, nil)
}

func Incr(name string, labels map[string]string) error {
	return metr.Incr(name, labels)
}

func Store(name string, value float64, labels map[string]string) error {
	return metr.Store(name, value, labels)
}

func RegisterGauge(name, desc string, labels []string) error {
	return metr.RegisterGauge(name, desc, labels)
}

func RegisterCounter(name, desc string, labels []string) error {
	return metr.RegisterCounter(name, desc, labels)
}
