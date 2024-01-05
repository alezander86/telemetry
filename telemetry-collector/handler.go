package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type MetricsHandler struct {
	Collector *Collector
}

func NewMetricsHandler(collector *Collector) *MetricsHandler {
	return &MetricsHandler{Collector: collector}
}

func (h *MetricsHandler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == http.MethodPost && request.Path == "/v1/submit" {
		err := h.Collector.Create(request.Body)
		if err != nil {
			log.Printf("Error processing metrics: %v\n", err)

			return events.APIGatewayProxyResponse{
				Body:       "Invalid request",
				StatusCode: http.StatusInternalServerError,
			}, nil
		}

		return events.APIGatewayProxyResponse{
			Body:       "ok",
			StatusCode: http.StatusOK,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "Invalid request",
		StatusCode: http.StatusBadRequest,
	}, nil
}
