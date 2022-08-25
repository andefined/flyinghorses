package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/andefined/flyinghorses/internal/services/monitor"
	"github.com/andefined/flyinghorses/pkg/config"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	// Create a new Config struct to store environment variables
	cfg := config.NewConfig()

	// Read config from env variables and panic on error,
	// as we can't continue
	err := envconfig.Process("", cfg)
	if err != nil {
		fmt.Printf("Failed to parse environment variables, exiting with error: %s\n", err.Error())
		os.Exit(1)
	}

	// Set MAXPROCS
	runtime.GOMAXPROCS(cfg.NumCPU)

	// Create the context to cancel at exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the gRPC Server
	if err := monitor.NewMonitorService(ctx, cfg); err != nil {
		fmt.Printf("CellMonitorServive failure, exiting with error: %s\n", err.Error())
		os.Exit(1)
	}
}
