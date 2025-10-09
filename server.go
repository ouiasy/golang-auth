package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ouiasy/golang-auth/api"
	"github.com/ouiasy/golang-auth/conf"
	"github.com/ouiasy/golang-auth/mailer"
	"github.com/ouiasy/golang-auth/repository"
)

func serve() error {
	initLogger()

	if os.Getenv("APP_ENV") == "local" {
		err := conf.LoadDotEnvDir()
		if err != nil {
			slog.Error("error while loading .env folder")
			return fmt.Errorf("error while loading .env folder: %v", err)
		}
	}

	config, err := conf.LoadConfigFromEnv()
	if err != nil {
		slog.Error("error while parsing env vars into Config struct")
		return err
	}

	repo, err := repository.NewRepository(config)
	if err != nil {
		slog.Error("error while establishing db connection")
		return err
	}

	emailClient := mailer.NewEmailClient(config)

	api := api.NewApi(config, repo, emailClient)

	server := &http.Server{
		Handler:           api.Handler,
		Addr:              api.Config.App.Host + ":" + api.Config.App.Port,
		ReadHeaderTimeout: 2 * time.Second, // against Slowloris attack
	}

	done := make(chan struct{})
	defer close(done)
	go func() {
		waitForTermination(done)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		server.Shutdown(ctx)
	}()

	fmt.Println("[+] running server at http://" + server.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error(err.Error())
		return err
	}

	return nil
}

// WaitForShutdown blocks until the system signals termination or done has a value
func waitForTermination(done <-chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-signals:
		slog.Info("Triggering shutdown from signal: " + sig.String())
	case <-done:
		slog.Info("Shutting down...")
	}
}

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//logger = logger.With("component", "api")
	slog.SetDefault(logger)
	if os.Getenv("APP_ENV") == "local" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}
