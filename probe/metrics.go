package probe

import "github.com/prometheus/client_golang/prometheus"

type collector interface {
	prometheus.Collector
}

type Metric interface {
	prometheus.Metric
}

type GaugeMetric interface {
	prometheus.Gauge
}

type CounterMetric interface {
	prometheus.Counter
}

type HistogramMetric interface {
	prometheus.Histogram
}
