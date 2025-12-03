package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"cron-scheduler/internal/config"
	"cron-scheduler/internal/scheduler"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	statePath := flag.String("state", "state.json", "path to state file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Println("config error:", err)
		return
	}

	state := scheduler.NewStateManager(*statePath)
	state.Load()

	s := scheduler.NewScheduler(cfg.Jobs, state)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s.Start(ctx)

	fmt.Println("GoScheduler running... press CTRL+C to stop")
	<-ctx.Done()

	s.Wait()
	fmt.Println("Shutdown complete")
}
