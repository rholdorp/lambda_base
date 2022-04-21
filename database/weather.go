package database

import (
	"context"
	"database/sql"
)

const createWeatherLogSQL = `
INSERT INTO 
	weather (in_temp, out_temp, humidity)
VALUES ($1, $2, $3)	
`

func CreateWeatherLog(ctx context.Context, db *sql.DB, in_temp string, out_temp string, humidity string) error {
	_, err := db.ExecContext(ctx, createWeatherLogSQL, in_temp, out_temp, humidity)
	return err
}

const GetWeatherSQL = `
SELECT * FROM weather
`

type WeatherLog struct {
	ID       int
	In_temp  string `json:"in_temp"`
	Out_temp string `json:"out_temp"`
	Humidity string `json:"humidity"`
}

func GetWeather(ctx context.Context, db *sql.DB) ([]*WeatherLog, error) {

	rows, err := db.QueryContext(ctx, GetWeatherSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	weather := make([]*WeatherLog, 0)

	for rows.Next() {
		weatherlog := &WeatherLog{}
		if err := rows.Scan(&weatherlog.ID, &weatherlog.In_temp,
			&weatherlog.Out_temp, &weatherlog.Humidity); err != nil {
			return nil, err
		}
		weather = append(weather, weatherlog)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return weather, nil
}
