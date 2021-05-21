package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb"
)

type sqlserverEngineRotator struct{}

func (s *sqlserverEngineRotator) openConnection(creds databaseCredentials) (*sql.DB, error) {
	host := creds.Host

	if creds.Port != "" {
		host += ":" + creds.Port
	}

	log.Printf("Connecting to SQL Server using %s@%s", creds.Username, host)

	query := url.Values{}
	query.Add("connection timeout", "30")
	query.Add("database", creds.DatabaseName)
	query.Add("dial timeout", "30")

	sqlServerDsn := &url.URL{
		Host:     host,
		Path:     creds.Instance,
		RawQuery: query.Encode(),
		Scheme:   "sqlserver",
		User:     url.UserPassword(creds.Username, creds.Password),
	}

	return sql.Open("sqlserver", sqlServerDsn.String())
}

func (s *sqlserverEngineRotator) Rotate(creds databaseCredentials, newPassword string) error {
	db, err := s.openConnection(creds)

	if err != nil {
		return fmt.Errorf("Could not connect to DB: %v", err)
	}

	defer db.Close()

	sql := fmt.Sprintf(
		"ALTER LOGIN [%s] WITH PASSWORD = '%s' OLD_PASSWORD = '%s'",
		creds.Username,
		newPassword,
		creds.Password,
	)

	if _, err = db.Exec(sql); err != nil {
		log.Printf("Error while changing password: %v", err)
		return err
	}

	return nil
}

func (s *sqlserverEngineRotator) Test(credentials databaseCredentials) error {
	db, err := s.openConnection(credentials)

	if err == nil {
		defer db.Close()

		err = db.Ping()
	}

	return err
}
