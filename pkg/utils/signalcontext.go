package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/italypaleale/revaulter/pkg/utils/applogger"
)

/*
This code is adapted from:
https://github.com/kubernetes-sigs/controller-runtime/blob/8499b67e316a03b260c73f92d0380de8cd2e97a1/pkg/manager/signals/signal.go
Copyright 2017 The Kubernetes Authors.
License: Apache2 (https://github.com/kubernetes-sigs/controller-runtime/blob/8499b67e316a03b260c73f92d0380de8cd2e97a1/LICENSE)
*/

var onlyOneSignalHandler = make(chan struct{})

// SignalContext returns a context that is canceled when the application receives an interrupt signal.
// A second signal forces an immediate shutdown.
func SignalContext(appLogger *applogger.Logger) context.Context {
	close(onlyOneSignalHandler) // Panics when called twice

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		appLogger.Raw().Info().Msg("Received interrupt signal. Shutting down…")
		cancel()

		<-sigCh
		appLogger.Raw().Fatal().Msg("Received a second interrupt signal. Forcing an immediate shutdown.")
	}()

	return ctx
}
