package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	db "github.com/praveenmahasena/server/internal/DB"
	"github.com/praveenmahasena/server/internal/server"
)
var S *server.Server

func Run() error {
	if dotEnvErr := godotenv.Load(".env"); dotEnvErr != nil {
		return fmt.Errorf("error during reading .env file %v", dotEnvErr)
	}
	libConnectionUri := db.LibDBUri()
	connection, connectionErr := db.NewConnection(libConnectionUri)
	if connectionErr != nil {
		return connectionErr
	}
	S = server.New(":42069", connection)
	errCh := make(chan error)
	sig := make(chan os.Signal, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // I do not like it here but I don't wanna have accidental context leak :)
	go func(e chan<- error) {
		e <- S.Start(ctx)
	}(errCh)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	// The errCh returns the val here if it has problem starting up the server
	case err := <-errCh:
		return err
	case <-sig:
		cancel()
		fmt.Fprintln(os.Stdout, "shutting down server...")
	}
	// The errCh returns the val here if it has problem closing up the server
	return <-errCh
}
