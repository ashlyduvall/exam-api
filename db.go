package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

type sqlTestRow struct {
	Value bool
}

func getConnectUrl() (string, error) {
	DB_PASS, found := os.LookupEnv("DB_PASS")
	if !found {
		return "", errors.New("Env var DB_PASS not set!")
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		GetConfig("DB_USER", "mysql"),
		DB_PASS,
		GetConfig("DB_HOST", "localhost"),
		GetConfig("DB_PORT", "3306"),
		GetConfig("DB_SCHEMA", "exam"),
	), nil
}

func TestDbConnectivity() error {
	url, err := getConnectUrl()

	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", url)

	if err != nil {
		return err
	}

	defer db.Close()
	var t sqlTestRow
	err = db.QueryRow("SELECT 1 AS test").Scan(&t.Value)

	if err != nil {
		return err
	}

	return nil
}
