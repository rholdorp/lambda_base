package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"lambda-base/database"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var logger *zap.Logger
var db *sql.DB

type DefaultResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GetWeatherResponse struct {
	Weather []*database.WeatherLog `json:"weather"`
}

func init() {
	l, _ := zap.NewProduction()
	logger = l

	dbConnection, err := database.GetConnection()
	if err != nil {
		logger.Error("error connection to database", zap.Error(err))
		panic(err)
	}

	dbConnection.Ping()
	if err != nil {
		logger.Error("error pinging the database", zap.Error(err))
		panic(err)
	}

	db = dbConnection
}

type Event struct {
	Name string `json:"name"`
}

func MyHandler(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	logger.Info("recieved event", zap.Any("method", event.HTTPMethod), zap.Any("path", event.Path), zap.Any("body", event.Body))

	if event.Path == "/migrate" {
		err := database.CreateWeatherTable(ctx, db)

		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  string(http.StatusInternalServerError),
				Message: "unable to create weather table",
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       string(body),
			}, nil
		}

		body, _ := json.Marshal(&DefaultResponse{
			Status:  string(http.StatusOK),
			Message: "table succesfully created",
		})
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil

	} else if event.Path == "/weatherlog" && event.HTTPMethod == http.MethodPost {
		// create weather
		weatherlog := &database.WeatherLog{}
		err := json.Unmarshal([]byte(event.Body), &weatherlog)

		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  string(http.StatusBadRequest),
				Message: err.Error(),
			})
			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       string(body),
			}, nil
		}

		err = database.CreateWeatherLog(ctx, db, weatherlog.In_temp, weatherlog.Out_temp, weatherlog.Humidity)
		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  string(http.StatusInternalServerError),
				Message: err.Error(),
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       string(body),
			}, nil
		}

		body, _ := json.Marshal(&DefaultResponse{
			Status:  string(http.StatusOK),
			Message: "Success!",
		})

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil

	} else if event.Path == "/weatherlog" && event.HTTPMethod == http.MethodGet {
		// get weather
		weather, err := database.GetWeather(ctx, db)

		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  string(http.StatusInternalServerError),
				Message: "unable to get weather",
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       string(body),
			}, nil
		}

		body, _ := json.Marshal(&GetWeatherResponse{
			Weather: weather,
		})

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil

	} else {
		body, _ := json.Marshal(&DefaultResponse{
			Status:  string(http.StatusOK),
			Message: "default path",
		})
		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}
	}
	return res, nil
}

func main() {
	lambda.Start(MyHandler)
}
