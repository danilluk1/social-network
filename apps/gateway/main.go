package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	execCtx, execCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	defer execCancel()

	go func() {
		<-execCtx.Done()
		logrus.Info("receivved graceful shutdown signal")
	}()

	if err := cmd.RootCommand().ExecuteContext(execCtx); err != nil {
		logrus.WithError(err).Fatal(err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
	defer shutdownCancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		api.WaitForCleanup(shutdownCtx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		observability.WaitForCleanup(shutdownCtx)
	}()

	cleanupDone := make(chan struct{})
	go func() {
		defer close(cleanupDone)
		wg.Wait()
	}()

	select {
	case <-shutdownCtx.Done():
		return

	case <-cleanupDone:
		return
	}
}
