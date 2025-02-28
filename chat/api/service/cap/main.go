package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Natnael-Alemayehu/chat/chat/foundation/logger"
)

func main() {
	var log *logger.Logger

	traceIDFn := func(ctx context.Context) string {
		return "" // TODO Later for message tracing
	}

	log = logger.New(os.Stdout, logger.LevelInfo, "CAP", traceIDFn)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	log.Info(ctx, "startup", "Status", "Started")
	defer log.Info(ctx, "Shuting Down", "Shutting Down")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	return nil
}
