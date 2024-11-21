package main

import (
	"context"
	_ "embed"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	_ "go.uber.org/automaxprocs"

	"github.com/buglloc/vanityd/internal/httpd"
	"github.com/buglloc/vanityd/internal/projects"
)

//go:embed projects.yaml
var embedProjects []byte

func main() {
	var sourcePath string
	var zone string
	var addr string
	flag.StringVar(&sourcePath, "source", "", "path to a source file")
	flag.StringVar(&zone, "zone", httpd.DefaultZone, "zone to use")
	flag.StringVar(&addr, "addr", httpd.DefaultAddr, "address to listen on")
	flag.Parse()

	var prj projects.Provider
	var err error
	if sourcePath != "" {
		prj, err = projects.NewFSProjects(sourcePath)
	} else {
		prj, err = projects.NewStaticProjects(embedProjects)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("could not create project provider")
	}

	http := httpd.NewHttpD(prj,
		httpd.WithAddr(addr),
		httpd.WithZone(zone),
	)

	errChan := make(chan error, 1)
	okChan := make(chan struct{}, 1)
	go func() {
		if err := http.ListenAndServe(); err != nil {
			errChan <- err
			return
		}

		okChan <- struct{}{}
	}()

	defer log.Info().Msg("stopped")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Info().Msg("shutting down gracefully by signal")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = http.Shutdown(ctx)
	case <-okChan:
		return
	case err := <-errChan:
		log.Fatal().Err(err).Msg("could not start server")
	}
}
