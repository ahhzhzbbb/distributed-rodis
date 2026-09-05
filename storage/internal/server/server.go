package server

import (
	"fmt"
	"net"
	"time"

	"rodis/storage/internal/engine"
)

type Server struct {
	Config
	ln  net.Listener
	kv  *engine.KeyValue
	sem chan struct{}
}

func NewServer(cfg Config) *Server {
	if len(cfg.Port) == 0 {
		cfg.Port = defaultPort
	}
	return &Server{
		Config: cfg,
		kv:     engine.NewKeyValue(),
		sem:    make(chan struct{}, cfg.MaxConnections),
	}
}

func (s *Server) Start() error {
	s.Banner()

	ln, err := net.Listen("tcp", s.Port)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.runActiveExpiration()

	go func() {
		for err := range s.loop() {
			fmt.Printf("ERROR: %s\n", err)
		}
	}()
	return nil
}

func (s *Server) loop() chan error {
	errStream := make(chan error)

	go func() {
		defer close(errStream)
		for {
			conn, err := s.ln.Accept()
			if err != nil {
				errStream <- err
				return
			}

			s.sem <- struct{}{}
			go func(conn net.Conn) {
				defer func() {
					<-s.sem
					_ = conn.Close()
				}()
				s.handleConnection(conn)
			}(conn)
		}
	}()

	return errStream
}

func (s *Server) runActiveExpiration() {
	for {
		time.Sleep(time.Duration(s.Expire.CycleIntervalMs) * time.Millisecond)
		s.kv.ActiveExpiration(s.Expire.SampleSize, s.Expire.ExpireThreshold, s.Expire.TimeBudgetMs)
	}
}
