package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/danilluk1/social-network/apps/auth/cmd"
	"github.com/danilluk1/social-network/apps/auth/internal/observability"
)

func main() {
	execCtx, execCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	defer execCancel()
	go func() {
		<-execCtx.Done()
		log.Println("Auth services shutted down with execCtx")
	}()

	// command is expected to obey the cancellation signal on execCtx and
	// block while it is running
	if err := cmd.RootCommand().ExecuteContext(execCtx); err != nil {
		log.Fatal(err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
	defer shutdownCancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		// wait for API servers to shut down gracefully
		// gapi.WaitForCleanup(shutdownCtx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// wait for profiler, metrics and trace exporters to shut down gracefully
		observability.WaitForCleanup(shutdownCtx)
	}()

	cleanupDone := make(chan struct{})
	go func() {
		defer close(cleanupDone)
		wg.Wait()
	}()

	select {
	case <-shutdownCtx.Done():
		// cleanup timed out
		return

	case <-cleanupDone:
		// cleanup finished before timing out
		return
	}
}
