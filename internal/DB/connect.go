package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	hostParam     = "HOST"
	portParam     = "PORT"
	userParam     = "USER"
	passwordParam = "PASSWORD"
	dbnameParam   = "DBNAME"
	sslmodeParam  = "SSLMODE"
)

func LibDBUri() string {
	host := os.Getenv(hostParam)
	port := os.Getenv(portParam)
	user := os.Getenv(userParam)
	password := os.Getenv(passwordParam)
	dbname := os.Getenv(dbnameParam)
	sslmode := os.Getenv(sslmodeParam)

	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v", host, port, user, password, dbname, sslmode)
}

func NewConnection(uri string) (*sql.DB ,error) {
	connection, connectionErr := sql.Open("postgres", uri)
	if connectionErr != nil {
		return nil, fmt.Errorf("error during connection to database %v", connectionErr)
	}
	if connectionPingErr := connection.Ping(); connectionPingErr != nil {
		return nil, fmt.Errorf("error during connection checkup %v", connectionPingErr)
	}
	return connection,nil
}

