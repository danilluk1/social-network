package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/danilluk1/social-network/apps/mailer/cmd"
)

func main() {
	execCtx, execCancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	defer execCancel()
	go func() {
		<-execCtx.Done()
		log.Println("Mailer service shutted down with execCtx")
	}()

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
