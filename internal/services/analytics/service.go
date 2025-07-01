package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/marcusziade/pomme/internal/client"
	"github.com/marcusziade/pomme/internal/models"
	"github.com/marcusziade/pomme/internal/services/cache"
)

// Service provides analytics functionality
type Service struct {
	client *client.Client
	cache  cache.Cache
	mu     sync.RWMutex
}

// NewService creates a new analytics service
func NewService(client *client.Client) *Service {
	return &Service{
		client: client,
		cache:  cache.NewMemoryCache(),
	}
}

// GetAppMetrics fetches performance and power metrics for an app
func (s *Service) GetAppMetrics(ctx context.Context, appID string) (*AnalyticsReport, error) {
	cacheKey := fmt.Sprintf("app_metrics_%s", appID)
	
	// Check cache
	if cached, err := s.cache.Get(cacheKey); err == nil {
		if report, ok := cached.(*AnalyticsReport); ok {
			return report, nil
		}
	}

	// Fetch metrics
	endpoint := fmt.Sprintf("/v1/apps/%s/perfPowerMetrics", appID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Include related resources
	q := req.URL.Query()
	q.Add("include", "diagnosticSignatures,diagnosticInsights")
	q.Add("limit", "200")
	req.URL.RawQuery = q.Encode()

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Data     []models.PerfPowerMetric      `json:"data"`
		Included []interface{}                 `json:"included"`
		Links    map[string]interface{}        `json:"links"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Process the response into a report
	report := s.processMetricsResponse(appID, response.Data, response.Included)
	
	// Cache the result
	s.cache.Set(cacheKey, report, 15*time.Minute)

	return report, nil
}

// GetBuildMetrics fetches metrics for a specific build
func (s *Service) GetBuildMetrics(ctx context.Context, buildID string) (*BuildMetrics, error) {
	endpoint := fmt.Sprintf("/v1/builds/%s/perfPowerMetrics", buildID)
	req, err := s.client.NewRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var metrics BuildMetrics
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &metrics, nil
}

// CompareVersions compares metrics between two app versions
func (s *Service) CompareVersions(ctx context.Context, appID, version1, version2 string) (*VersionComparison, error) {
	// This would require fetching metrics for both versions and computing differences
	// Implementation would depend on the exact API structure
	return nil, fmt.Errorf("not implemented")
}

// processMetricsResponse converts raw API response into structured report
func (s *Service) processMetricsResponse(appID string, metrics []models.PerfPowerMetric, included []interface{}) *AnalyticsReport {
	report := &AnalyticsReport{
		AppID:      appID,
		FetchedAt:  time.Now(),
		Categories: make(map[string]*models.MetricCategory),
	}

	// Group metrics by type
	launchMetrics := &models.MetricCategory{
		Identifier:  "launch_time",
		DisplayName: "Launch Time",
		Metrics:     []models.Metric{},
	}

	memoryMetrics := &models.MetricCategory{
		Identifier:  "memory",
		DisplayName: "Memory Usage",
		Metrics:     []models.Metric{},
	}

	batteryMetrics := &models.MetricCategory{
		Identifier:  "battery",
		DisplayName: "Battery Usage",
		Metrics:     []models.Metric{},
	}

	performanceMetrics := &models.MetricCategory{
		Identifier:  "performance",
		DisplayName: "Performance",
		Metrics:     []models.Metric{},
	}

	// Process each metric
	for _, metric := range metrics {
		switch metric.Attributes.MetricType {
		case "LAUNCH_TIME":
			launchMetrics.Metrics = append(launchMetrics.Metrics, s.createMetric(metric))
		case "MEMORY_USAGE":
			memoryMetrics.Metrics = append(memoryMetrics.Metrics, s.createMetric(metric))
		case "BATTERY_USAGE":
			batteryMetrics.Metrics = append(batteryMetrics.Metrics, s.createMetric(metric))
		default:
			performanceMetrics.Metrics = append(performanceMetrics.Metrics, s.createMetric(metric))
		}
	}

	// Add categories to report
	if len(launchMetrics.Metrics) > 0 {
		report.Categories["launch_time"] = launchMetrics
	}
	if len(memoryMetrics.Metrics) > 0 {
		report.Categories["memory"] = memoryMetrics
	}
	if len(batteryMetrics.Metrics) > 0 {
		report.Categories["battery"] = batteryMetrics
	}
	if len(performanceMetrics.Metrics) > 0 {
		report.Categories["performance"] = performanceMetrics
	}

	// Process insights from included data
	for _, item := range included {
		if data, ok := item.(map[string]interface{}); ok {
			if data["type"] == "diagnosticInsights" {
				// Process insights
				report.Insights = append(report.Insights, s.processInsight(data))
			}
		}
	}

	return report
}

// createMetric creates a metric from raw data
func (s *Service) createMetric(raw models.PerfPowerMetric) models.Metric {
	// This is a simplified version - real implementation would parse actual metric data
	return models.Metric{
		Identifier:  raw.ID,
		DisplayName: s.formatMetricName(raw.Attributes.MetricType),
		Unit:        s.getMetricUnit(raw.Attributes.MetricType),
		Percentiles: models.Percentiles{
			P50: models.MetricValue{Value: 0},
			P90: models.MetricValue{Value: 0},
			P95: models.MetricValue{Value: 0},
		},
	}
}

// formatMetricName converts metric type to display name
func (s *Service) formatMetricName(metricType string) string {
	switch metricType {
	case "LAUNCH_TIME":
		return "Time to First Frame"
	case "MEMORY_USAGE":
		return "Peak Memory"
	case "BATTERY_USAGE":
		return "Battery Drain"
	case "HANG_RATE":
		return "Hang Rate"
	case "DISK_WRITES":
		return "Disk Writes"
	default:
		return metricType
	}
}

// getMetricUnit returns the unit for a metric type
func (s *Service) getMetricUnit(metricType string) string {
	switch metricType {
	case "LAUNCH_TIME":
		return "ms"
	case "MEMORY_USAGE":
		return "MB"
	case "BATTERY_USAGE":
		return "%/hr"
	case "HANG_RATE":
		return "per hour"
	case "DISK_WRITES":
		return "MB"
	default:
		return ""
	}
}

// processInsight processes diagnostic insight data
func (s *Service) processInsight(data map[string]interface{}) models.DiagnosticInsight {
	// Simplified processing - real implementation would properly parse the data
	return models.DiagnosticInsight{}
}

// AnalyticsReport represents a complete analytics report
type AnalyticsReport struct {
	AppID      string                          `json:"appId"`
	FetchedAt  time.Time                       `json:"fetchedAt"`
	Categories map[string]*models.MetricCategory `json:"categories"`
	Insights   []models.DiagnosticInsight      `json:"insights"`
	Versions   []models.AppVersion             `json:"versions"`
}

// BuildMetrics represents metrics for a specific build
type BuildMetrics struct {
	BuildID    string                          `json:"buildId"`
	Version    string                          `json:"version"`
	Categories map[string]*models.MetricCategory `json:"categories"`
}

// VersionComparison represents a comparison between versions
type VersionComparison struct {
	Version1    string                   `json:"version1"`
	Version2    string                   `json:"version2"`
	Changes     []models.MetricSummary   `json:"changes"`
	Improvements int                     `json:"improvements"`
	Regressions int                     `json:"regressions"`
}