package probe

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type server struct {
	config     *ServerConfig
	registry   *prometheus.Registry
	httpServer *http.Server
	closed     chan error
}

func newServer(config *ServerConfig) server {
	registry := prometheus.NewRegistry()
	return server{
		config:   config,
		closed:   make(chan error, 1),
		registry: registry,
		httpServer: &http.Server{
			Addr:    config.Address,
			Handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		},
	}
}

func (s *server) register(collector collector) {
	s.registry.Register(collector)
}

func (s *server) start() {
	go func() {
		s.closed <- s.httpServer.ListenAndServe()
	}()
}

func (s *server) stop() error {
	return s.httpServer.Close()
}
