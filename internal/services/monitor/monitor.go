package monitor

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andefined/flyinghorses/internal/services/cell"
	"github.com/andefined/flyinghorses/pkg/config"
	"github.com/andefined/flyinghorses/pkg/logger"
	"github.com/go-redis/redis/v8"
)

// NewMonitorService
func NewMonitorService(ctx context.Context, cfg *config.Config) error {
	// ** LOGGER
	// Create a reusable zap logger
	log := logger.NewLogger(cfg.Env, cfg.Log.Level, cfg.Log.Path)
	log.Info("Subscribe to srsLTE->cell_measurement messages")

	// ** TERMINATION
	// Listen for os signals
	osSignals := make(chan os.Signal, 1)
	defer close(osSignals)
	// Listen for manual termination
	signal.Notify(osSignals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // We connect to host redis, thats what the hostname of the redis service is set to in the docker-compose
	})
	// Ping the Redis server and check if any errors occured
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}

	// Create the error channel
	errorChannel := make(chan error, 1)
	defer close(errorChannel)

	topicChannel := redisClient.Subscribe(ctx, "cell")
	cellMeasurementService := cell.NewCellMeasurementSerice(log, ctx, topicChannel, errorChannel)
	go cellMeasurementService.Consume()

	select {
	case err := <-errorChannel:
		log.Errorf("srsLTE->cell_measurement error: %s", err.Error())
		return err
	case signal := <-osSignals:
		log.Fatalf("srsLTE->cell_measurement shutdown signal: %s", signal)
	}

	return nil
}
