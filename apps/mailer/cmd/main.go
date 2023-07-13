package main

import (
	"context"
	"net/http"

	"github.com/danilluk1/social-network/libs/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if cfg.AppEnv != "development" {
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe("0.0.0.0:3000", nil)
	}

	var logger *zap.Logger

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	zap.ReplaceGlobals(logger)

}
