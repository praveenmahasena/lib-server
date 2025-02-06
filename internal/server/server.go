package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

type port string

type Server struct {
	port
	*sql.DB
}

func New(port port, db *sql.DB) *Server {
	return &Server{
		port,
		db,
	}
}

func (s Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()


	mux.HandleFunc("/idx",s.HandleIdx)

	ser := http.Server{
		Addr:    string(s.port),
		Handler: mux,
	}

	errCh := make(chan error)
	go func(e chan<- error) {
		defer close(e)
		e <- fmt.Errorf("error during spwaning the server %v", ser.ListenAndServe())
	}(errCh)

	cttx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Minute)
	defer cancel()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		cancel()
		if e := ser.Shutdown(cttx); e != nil {
			return fmt.Errorf("error during shutting down the server %v", e)
		}
	}
	return nil
}

func (c *Server) GetCon() (*sql.DB, error) {
	if err := c.Ping(); err != nil {
		return nil, fmt.Errorf("error during checking up connection to DB %v", err)
	}
	return c.DB, nil
}
