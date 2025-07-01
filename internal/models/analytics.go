package models

import "time"

// PerfPowerMetric represents performance and power metrics for an app
type PerfPowerMetric struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes PerfPowerMetricAttrs   `json:"attributes"`
	Links      map[string]interface{} `json:"links"`
}

// PerfPowerMetricAttrs contains the metric attributes
type PerfPowerMetricAttrs struct {
	DeviceType  string `json:"deviceType"`
	MetricType  string `json:"metricType"`
	Platform    string `json:"platform"`
}

// DiagnosticSignature represents a diagnostic signature
type DiagnosticSignature struct {
	ID         string                       `json:"id"`
	Type       string                       `json:"type"`
	Attributes DiagnosticSignatureAttrs     `json:"attributes"`
	Links      map[string]interface{}       `json:"links"`
}

// DiagnosticSignatureAttrs contains diagnostic signature attributes
type DiagnosticSignatureAttrs struct {
	DiagnosticType string  `json:"diagnosticType"`
	Signature      string  `json:"signature"`
	Weight         float64 `json:"weight"`
}

// DiagnosticInsight represents insights about app performance
type DiagnosticInsight struct {
	ID         string                     `json:"id"`
	Type       string                     `json:"type"`
	Attributes DiagnosticInsightAttrs     `json:"attributes"`
}

// DiagnosticInsightAttrs contains insight attributes
type DiagnosticInsightAttrs struct {
	InsightType string  `json:"insightType"`
	Direction   string  `json:"direction"`  // REGRESSION or IMPROVEMENT
	Value       float64 `json:"referenceVersionValue"`
	Change      float64 `json:"latestVersionValue"`
}

// MetricCategory represents a category of metrics
type MetricCategory struct {
	Identifier  string   `json:"identifier"`
	DisplayName string   `json:"displayName"`
	Metrics     []Metric `json:"metrics"`
}

// Metric represents a specific performance metric
type Metric struct {
	Identifier  string      `json:"identifier"`
	DisplayName string      `json:"displayName"`
	Unit        string      `json:"unit"`
	Percentiles Percentiles `json:"percentiles"`
	Goal        *MetricGoal `json:"goal,omitempty"`
}

// Percentiles contains percentile values for a metric
type Percentiles struct {
	P50 MetricValue `json:"p50"`
	P90 MetricValue `json:"p90"`
	P95 MetricValue `json:"p95"`
}

// MetricValue represents a metric value with device breakdown
type MetricValue struct {
	Value           float64                   `json:"value"`
	DeviceBreakdown map[string]DeviceMetric   `json:"deviceBreakdown"`
}

// DeviceMetric represents metrics for a specific device
type DeviceMetric struct {
	Value      float64 `json:"value"`
	DeviceType string  `json:"deviceType"`
	Count      int     `json:"count"`
}

// MetricGoal represents a performance goal
type MetricGoal struct {
	Value       float64 `json:"value"`
	Percentile  string  `json:"percentile"`
	Description string  `json:"description"`
}

// AppVersion represents version-specific metric data
type AppVersion struct {
	Version         string          `json:"version"`
	BuildNumber     string          `json:"buildNumber"`
	ReleaseDate     time.Time       `json:"releaseDate"`
	MetricSummaries []MetricSummary `json:"metricSummaries"`
}

// MetricSummary provides a summary of a metric across versions
type MetricSummary struct {
	MetricID        string  `json:"metricId"`
	CurrentValue    float64 `json:"currentValue"`
	PreviousValue   float64 `json:"previousValue"`
	PercentChange   float64 `json:"percentChange"`
	Status          string  `json:"status"` // IMPROVED, REGRESSED, STABLE
}