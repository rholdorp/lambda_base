package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("host")
	password = os.Getenv("password")
	port     = 5432
	user     = "postgres"
	database = "postgres"
)

func GetConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, database)
	return sql.Open("postgres", psqlInfo)
}

const createWeatherTableSQL = `
CREATE TABLE weather (
	id serial,
	in_temp varchar,
	out_temp varchar,
	humidity varchar
);
`

func CreateWeatherTable(ctx context.Context, db *sql.DB) error {

	_, err := db.ExecContext(ctx, createWeatherTableSQL)
	return err
}
