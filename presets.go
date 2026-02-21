// Copyright 2026 Xavier Portilla Edo
// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package opentelemetry

import (
	"log/slog"
	"sync"
	"time"
)

// PresetConfig provides pre-configured setups for common scenarios.
type PresetConfig struct {
	// The base config to extend
	Config

	// Preset type
	Type PresetType
}

// PresetType defines common OpenTelemetry setup presets.
type PresetType string

const (
	// PresetJaeger configures for Jaeger tracing
	PresetJaeger PresetType = "jaeger"

	// PresetPrometheus configures for Prometheus metrics
	PresetPrometheus PresetType = "prometheus"

	// PresetConsole configures for console output (development)
	PresetConsole PresetType = "console"

	// PresetOTLP configures for standard OTLP (default)
	PresetOTLP PresetType = "otlp"
)

// NewWithPreset creates a new OpenTelemetry plugin with a preset configuration.
func NewWithPreset(preset PresetType, customConfig ...Config) *OpenTelemetry {
	config := createPresetConfig(preset)

	// Apply custom config if provided
	if len(customConfig) > 0 {
		mergeConfig(&config, customConfig[0])
	}

	// Create the plugin with preset type information and properly initialize all fields
	config.setDefaults()
	return &OpenTelemetry{
		config:     config,
		presetType: &preset,
		serverWg:   &sync.WaitGroup{},
	}
}

// createPresetConfig creates a config based on the preset type.
func createPresetConfig(preset PresetType) Config {
	switch preset {
	case PresetJaeger:
		return Config{
			OTLPEndpoint:   "http://localhost:14268/api/traces", // Jaeger OTLP HTTP endpoint
			OTLPUseHTTP:    true,
			ServiceName:    "genkit-service",
			MetricInterval: 30 * time.Second,
			LogLevel:       slog.LevelInfo,
		}

	case PresetPrometheus:
		return Config{
			ServiceName:              "genkit-service",
			MetricInterval:           15 * time.Second, // Prometheus scrapes frequently
			LogLevel:                 slog.LevelInfo,
			EnablePrometheusEndpoint: true,
			PrometheusPort:           9090,
			EnablePrometheusExporter: true, // Force Prometheus setup
			// MetricExporter will be set to Prometheus in the setup
		}

	case PresetConsole:
		return Config{
			ServiceName:    "genkit-service",
			MetricInterval: 10 * time.Second,
			LogLevel:       slog.LevelDebug,
			ForceExport:    true, // Always export in console mode
			MetricExporter: createStdoutMetricExporter(),
		}

	case PresetOTLP:
		fallthrough
	default:
		return Config{
			OTLPEndpoint:   "http://localhost:4317",
			OTLPUseHTTP:    false, // Use gRPC by default
			ServiceName:    "genkit-service",
			MetricInterval: 60 * time.Second,
			LogLevel:       slog.LevelInfo,
		}
	}
}

// mergeConfig merges custom config into the base config.
func mergeConfig(base *Config, custom Config) {
	if custom.ForceExport {
		base.ForceExport = custom.ForceExport
	}
	if custom.MetricInterval != 0 {
		base.MetricInterval = custom.MetricInterval
	}
	if custom.LogLevel != nil {
		base.LogLevel = custom.LogLevel
	}
	if custom.TraceExporter != nil {
		base.TraceExporter = custom.TraceExporter
	}
	if custom.MetricExporter != nil {
		base.MetricExporter = custom.MetricExporter
	}
	if custom.LogHandler != nil {
		base.LogHandler = custom.LogHandler
	}
	if custom.OTLPEndpoint != "" {
		base.OTLPEndpoint = custom.OTLPEndpoint
	}
	if custom.OTLPUseHTTP {
		base.OTLPUseHTTP = custom.OTLPUseHTTP
	}
	if custom.OTLPHeaders != nil {
		if base.OTLPHeaders == nil {
			base.OTLPHeaders = make(map[string]string)
		}
		for k, v := range custom.OTLPHeaders {
			base.OTLPHeaders[k] = v
		}
	}
	if custom.ServiceName != "" {
		base.ServiceName = custom.ServiceName
	}
	if custom.ServiceVersion != "" {
		base.ServiceVersion = custom.ServiceVersion
	}
	if custom.ResourceAttributes != nil {
		if base.ResourceAttributes == nil {
			base.ResourceAttributes = make(map[string]string)
		}
		for k, v := range custom.ResourceAttributes {
			base.ResourceAttributes[k] = v
		}
	}
	if custom.EnablePrometheusEndpoint {
		base.EnablePrometheusEndpoint = custom.EnablePrometheusEndpoint
	}
	if custom.PrometheusPort != 0 {
		base.PrometheusPort = custom.PrometheusPort
	}
	if custom.EnablePrometheusExporter {
		base.EnablePrometheusExporter = custom.EnablePrometheusExporter
	}
}
