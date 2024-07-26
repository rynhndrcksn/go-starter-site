package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serve sets up a new http.Server and calls .ListenAndServe on it.
func (app *application) serve() error {
	// Initialize a tls.Config struct to hold any non-default TLS settings.
	// In this case, change it to only use elliptic curves with assembly implementations that Go supports.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	// Initialize HTTP server using some sensible timeout settings.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.sessionManager.LoadAndSave(app.routes()),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		TLSConfig:    tlsConfig,
		WriteTimeout: 10 * time.Second,
	}

	// Create a shutdownError channel used to receive any errors returned by the graceful Shutdown() function.
	shutdownError := make(chan error)

	// Start a background goroutine.
	go func() {
		// Create a quit channel which carries os.Signal values.
		quit := make(chan os.Signal, 1)

		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and relay them to the quit channel.
		// Any other signals will not be caught by signal.Notify() and will retain their default behavior.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Read the signal from the quit channel.
		// This code will block until a signal is received.
		s := <-quit

		// Log a message to say that the signal has been caught.
		app.logger.Info("shutting down server", slog.String("signal", s.String()))

		// Create a context with a 30-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Call Shutdown() on the server, passing in the context that was just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an error (which may happen because of a
		// problem closing the listeners, or because the shutdown didn't complete before the 30-second context deadline is hit).
		// Relay this return value to the shutdownError channel.
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks", slog.String("addr", srv.Addr))

		// Call Wait() to block until the WaitGroup counter is zero, blocking until the background goroutines have finished.
		// Then return nil on the shutdownError channel to indicate that the shutdown was completed without any issues.
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	// Calling Shutdown() on the server will cause ListenAndServe() to immediately return a http.ErrServerClosed error.
	// So if this error is seen, it is actually a good thing and a sign that the graceful shutdown has started.
	// So check specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, wait to receive the return value from Shutdown() on the shutdownError channel.
	// If return value is an error, there was a problem with the graceful shutdown, and return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("server stopped", slog.String("addr", srv.Addr))
	return nil
}
