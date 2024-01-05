package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/opensearch-project/opensearch-go/v2"
)

const (
	metricsIndex = "metrics-index"
)

type CodebaseMetrics struct {
	Lang       string `json:"lang"`
	Framework  string `json:"framework"`
	BuildTool  string `json:"buildTool"`
	Strategy   string `json:"strategy"`
	Type       string `json:"type"`
	Versioning string `json:"versioning"`
}

type CdPipelineMetrics struct {
	DeploymentType string `json:"deploymentType"`
	NumberOfStages int    `json:"numberOfStages"`
}

type PlatformMetrics struct {
	CodebaseMetrics   []CodebaseMetrics   `json:"codebaseMetrics"`
	CdPipelineMetrics []CdPipelineMetrics `json:"cdPipelineMetrics"`
	GitProviders      []string            `json:"gitProviders"`
	JiraEnabled       bool                `json:"jiraEnabled"`
	RegistryType      string              `json:"registryType"`
	Version           string              `json:"version"`
}

type PlatformMetricsDocument struct {
	Timestamp       time.Time        `json:"@timestamp"`
	PlatformMetrics *PlatformMetrics `json:"platformMetrics"`
}

type Collector struct {
	opensearch *opensearch.Client
}

func NewCollector(opensearch *opensearch.Client) *Collector {
	return &Collector{opensearch: opensearch}
}

// Create pushes a new document to the specified index
func (c Collector) Create(request string) error {
	var metrics PlatformMetricsDocument

	err := json.Unmarshal([]byte(request), &metrics)
	if err != nil {
		return fmt.Errorf("failed to unmarshal request body: %v", err)
	}

	if metrics.PlatformMetrics == nil {
		return fmt.Errorf("invalid request body")
	}

	metrics.Timestamp = time.Now()

	doc, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %v", err)
	}

	log.Println("Metrics: ", metrics)

	response, err := c.opensearch.Index(metricsIndex, bytes.NewReader(doc))
	if err != nil {
		return fmt.Errorf("failed to add document to index: %v", err)
	}

	if response.IsError() {
		return fmt.Errorf("failed to add document to index: %s", response)
	}

	return nil
}
