package probe

import "github.com/prometheus/client_golang/prometheus"

type _probe struct {
	config  *Config
	server  server
	metrics map[string]Metric
}

var probe *_probe

// Start creates a probe which serves metrics via HTTP server, if it is enabled by config and not already started.
// Start is not thread-safe, caller should use due guards if necessary
func Start() {
	if probe != nil {
		return
	}
	config := GetConfig()
	if config.Enabled {
		probe = &_probe{
			config:  config,
			server:  newServer(config.ServerConfig),
			metrics: make(map[string]Metric),
		}
		probe.server.start()
	}
}

// Stop stops the HTTP server. Stop is not thread-safe, caller should use due guards if necessary
func Stop() {
	if probe != nil {
		probe.server.stop()
		probe = nil
	}
}

// Get gets the metric for the given metricId. Caller can use the returned metric by casting it to one of appropriate
// metric types; probe.GaugeMetric, probe.CounterMetric, probe.HistogramMetric
// - On success, returns (metric, true)
// - If probe is disabled or not started returns (nil, false)
func Get(metricId string) (Metric, bool) {
	if Enabled() {
		m, ok := probe.metrics[metricId]
		return m, ok
	}
	return nil, false
}

// Set creates a metric with given metricConfig if probe is enabled and started. If config file contains a metric with
// same id as metricConfig.Id, metricConfig values are overwritten by config file. Caller can use the returned metric by
// casting it to one of appropriate metric types; probe.GaugeMetric, probe.CounterMetric, probe.HistogramMetric
//  - On success, returns (metric, true)
//  - If probe is disabled or not started, returns (nil, false)
//  - If there is a metric created with same id previously, returns the previously created metric.
func Set(metricConfig MetricConfig) (Metric, bool) {
	var m Metric
	if Enabled() {

		// return if it is already defined
		if m, ok := Get(metricConfig.Name); ok {
			return m, ok
		}

		// check config and overwrite given metricConfig with config if id matches
		for _, mc := range *probe.config.MetricConfigs {
			if mc.Id == metricConfig.Id {
				if mc.Name != "" {
					metricConfig.Name = mc.Name
				}
				if mc.Help != "" {
					metricConfig.Help = mc.Help
				}
				if mc.LabelValues != nil && mc.LabelNames != nil {
					if len(mc.LabelValues) > 0 && len(mc.LabelNames) > 0 {
						metricConfig.LabelNames = mc.LabelNames
						metricConfig.LabelValues = mc.LabelValues
					}
				}
				if mc.Buckets != nil {
					metricConfig.Buckets = mc.Buckets
				}
			}
		}

		switch metricConfig.Type {
		case "Gauge":
			c := prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: probe.server.config.Namespace,
				Subsystem: probe.server.config.Subsystem,
				Name:      metricConfig.Name,
				Help:      metricConfig.Help,
			}, metricConfig.LabelNames)
			probe.server.register(c)
			m = c.WithLabelValues(metricConfig.LabelValues...)
			break
		case "Counter":
			c := prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: probe.server.config.Namespace,
				Subsystem: probe.server.config.Subsystem,
				Name:      metricConfig.Name,
				Help:      metricConfig.Help,
			}, metricConfig.LabelNames)
			probe.server.register(c)
			m = c.WithLabelValues(metricConfig.LabelValues...)
			break
		case "Histogram":
			c := prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: probe.server.config.Namespace,
				Subsystem: probe.server.config.Subsystem,
				Name:      metricConfig.Name,
				Help:      metricConfig.Help,
				Buckets:   nil,
			}, metricConfig.LabelNames)
			probe.server.register(c)
			m = c.WithLabelValues(metricConfig.LabelValues...).(HistogramMetric)
			break
		default:
		}
		probe.metrics[metricConfig.Id] = m
		return m, true
	}
	return nil, false
}

// Enabled returns true if probe is started
func Enabled() bool {
	return probe != nil
}
