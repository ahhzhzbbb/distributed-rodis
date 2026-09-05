package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"rodis/gateway/internal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	config := internal.NewConfig("0.0.0.0", ":6349", "8.8.8.8")
	logger := internal.NewLogger(zerolog.InfoLevel)

	proxy := internal.NewProxy(
		config,
		logger,
	)

	serverAddr := net.JoinHostPort(config.IpAddr, config.Port)

	http2Server := http.Server{
		Addr:    serverAddr,
		Handler: proxy,
	}

	http2Server.Protocols = new(http.Protocols)
	http2Server.Protocols.SetUnencryptedHTTP2(true)

	go func() {
		fmt.Fprintf(w, "listening on %s\n", http2Server.Addr)
		if err := http2Server.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel = context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()

		if err := http2Server.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()

	return nil
}

func main() {
	ctx := context.Background()

	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Println("Your gateway is shutting down...")
		os.Exit(1)
	}
}
