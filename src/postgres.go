package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/jackc/pgx/v4/stdlib" // Import needed for pgx driver
)

type postgresEngineRotator struct{}

func (p *postgresEngineRotator) openConnection(creds databaseCredentials) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?connect_timeout=5&sslmode=require",
		creds.Username, url.PathEscape(creds.Password), creds.Host, creds.Port, creds.DatabaseName,
	)

	return sql.Open("pgx", dsn)
}

func (p *postgresEngineRotator) Rotate(creds databaseCredentials, newPassword string) error {
	db, err := p.openConnection(creds)

	if err != nil {
		return fmt.Errorf("Could not connect to DB: %v", err)
	}

	defer db.Close()

	if _, err = db.Exec(fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s'", creds.Username, newPassword)); err != nil {
		log.Printf("Error while changing password: %v", err)
		return err
	}

	return nil
}

func (p *postgresEngineRotator) Test(credentials databaseCredentials) error {
	db, err := p.openConnection(credentials)

	if err == nil {
		defer db.Close()

		err = db.Ping()
	}

	return err
}
