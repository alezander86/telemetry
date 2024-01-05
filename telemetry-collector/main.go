package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/opensearch-project/opensearch-go/v2"
)

func main() {
	h := NewMetricsHandler(NewCollector(getElasticClient()))

	lambda.StartWithOptions(h.HandleRequest)
}

func getElasticClient() *opensearch.Client {
	elasticsearchURL := os.Getenv("ELASTIC_URL")
	if elasticsearchURL == "" {
		log.Fatal("ELASTIC_URL environment variable not set")
	}

	username := os.Getenv("ELASTIC_USERNAME")
	if username == "" {
		log.Fatal("ELASTIC_USERNAME environment variable not set")
	}

	password := os.Getenv("ELASTIC_PASSWORD")
	if password == "" {
		log.Fatal("ELASTIC_PASSWORD environment variable not set")
	}

	cfg := opensearch.Config{
		Addresses: []string{
			elasticsearchURL,
		},
		Username: username,
		Password: password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client, err := opensearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := client.Ping()
	if err != nil {
		log.Fatalf("Error pinging the client: %s", err)
	}

	if res.IsError() {
		log.Fatalf("Error pinging the client: %s", res)
	}

	return client
}
