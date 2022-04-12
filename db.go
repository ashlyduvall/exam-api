package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type sqlTestRow struct {
	Value bool
}

var DB *sql.DB

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

func ConnectAndTestDB() error {
	url, err := getConnectUrl()

	if err != nil {
		return err
	}

	pool, err := sql.Open("mysql", url)
	DB = pool

	if err != nil {
		return err
	}

	var t sqlTestRow
	err = DB.QueryRow("SELECT 1 AS test").Scan(&t.Value)

	if err != nil {
		return err
	}

	return nil
}
