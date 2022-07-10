package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/toqns/toqns/business/logo"
	"github.com/toqns/toqns/business/node"
	"github.com/toqns/toqns/foundation/logger"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	// Construct the application logger.
	log, err := logger.New("TOQNS", build == "develop")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// =========================================================================
	// Configuration

	// This is all the configuration for the application and the default values.
	// Configuration values will be passed through the application as individual
	// values.
	cfg := struct {
		conf.Version
		P2P struct {
			Address         string        `conf:"default:0.0.0.0"`
			Port            int           `conf:"default:3000"`
			Protocol        string        `conf:"default:udp"`
			NodeKeyFile     string        `conf:"default:./.node/node.key"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	// Parse will set the defaults and then look for any overriding values
	// in environment variables and command line flags.
	const prefix = "TOQNS"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// =========================================================================
	// App Starting

	logo.Print()
	fmt.Print("\n")

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	// Display the current configuration to the logs.
	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// =========================================================================
	// Service Start/Stop Support

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// =========================================================================
	// Start Public Service

	log.Infow("startup", "status", "initializing p2p support")

	n, err := node.New(log, node.NodeConfig{
		Address:     cfg.P2P.Address,
		Port:        cfg.P2P.Port,
		Protocol:    cfg.P2P.Protocol,
		NodeKeyFile: cfg.P2P.NodeKeyFile,
	})
	if err != nil {
		return fmt.Errorf("setting up p2p node: %w", err)
	}

	// Start the service listening for api requests.
	go func() {
		log.Infow("startup", "status", "p2p network started", "node address", n.Address.String())

		if err := n.ListenAndServe(); err != nil {
			serverErrors <- err
		}
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.P2P.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		log.Infow("shutdown", "status", "shutdown p2p network started")
		if err := n.Shutdown(ctx); err != nil {
			return fmt.Errorf("could not stop p2p network gracefully: %w", err)
		}
	}

	return nil
}
