package main

import (
	"errors"
	"fmt"
	"os"

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
	log, err := logger.New("TOQNS")
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
			Address     string `conf:"default:0.0.0.0"`
			Port        int    `conf:"default:3000"`
			Protocol    string `conf:"default:udp"`
			NodeKeyFile string `conf:"default:./.node/node.key"`
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

	n, err := node.New(log, node.NodeConfig{
		Address:     cfg.P2P.Address,
		Port:        cfg.P2P.Port,
		Protocol:    cfg.P2P.Protocol,
		NodeKeyFile: cfg.P2P.NodeKeyFile,
	})
	if err != nil {
		return fmt.Errorf("setting up p2p node: %w", err)
	}

	log.Infow("startup", "node address", n.Address.String())

	if err := n.ListenAndServe(); err != nil {
		return fmt.Errorf("starting listener: %w", err)
	}

	return nil
}
