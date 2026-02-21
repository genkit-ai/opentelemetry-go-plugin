// Copyright 2026 Xavier Portilla Edo
// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

// Package opentelemetry provides an OpenTelemetry plugin for Genkit Go.
// This plugin configures OpenTelemetry exporters for traces, metrics, and logs
// with sensible defaults while allowing customization of exporters.
package opentelemetry

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/firebase/genkit/go/core/api"
	"github.com/firebase/genkit/go/core/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Config configures the OpenTelemetry plugin.

const providerID = "opentelemetry"

type Config struct {
	// Export even in the dev environment.
	ForceExport bool

	// The interval for exporting metric data.
	// The default is 60 seconds.
	MetricInterval time.Duration

	// The minimum level at which logs will be written.
	// Defaults to slog.LevelInfo.
	LogLevel slog.Leveler

	// Custom trace span exporter. If nil, uses the default OTLP exporter.
	TraceExporter trace.SpanExporter

	// Custom metric exporter. If nil, uses the default OTLP exporter.
	MetricExporter metric.Exporter

	// Custom log handler. If nil, uses the default structured log handler.
	LogHandler slog.Handler

	// OTLP endpoint for traces and metrics. Defaults to "localhost:4317"
	// For gRPC (default), use format "host:port" (e.g., "localhost:4317")
	// For HTTP, use full URL format "http://host:port" or "https://host:port"
	OTLPEndpoint string

	// Whether to use HTTP instead of gRPC for OTLP. Defaults to false (gRPC).
	OTLPUseHTTP bool

	// Headers to include in OTLP requests.
	OTLPHeaders map[string]string

	// Service name for telemetry data. Defaults to "genkit-service".
	ServiceName string

	// Service version for telemetry data.
	ServiceVersion string

	// Additional resource attributes to include in telemetry data.
	ResourceAttributes map[string]string

	// Enable Prometheus metrics HTTP endpoint at /metrics. Defaults to false.
	EnablePrometheusEndpoint bool

	// Port for the Prometheus metrics HTTP server. Defaults to 9090.
	PrometheusPort int

	// Force enable Prometheus metrics setup regardless of preset type. Defaults to false.
	EnablePrometheusExporter bool
}

// setDefaults sets default values for the config.
func (c *Config) setDefaults() {
	if c.MetricInterval == 0 {
		c.MetricInterval = 60 * time.Second
	}
	if c.LogLevel == nil {
		c.LogLevel = slog.LevelInfo
	}
	if c.OTLPEndpoint == "" {
		c.OTLPEndpoint = "localhost:4317"
	}
	if c.ServiceName == "" {
		c.ServiceName = "genkit-service"
	}
	if c.ResourceAttributes == nil {
		c.ResourceAttributes = make(map[string]string)
	}
	if c.PrometheusPort == 0 {
		c.PrometheusPort = 9090
	}
}

// OpenTelemetry represents the OpenTelemetry plugin.
type OpenTelemetry struct {
	config       Config
	presetType   *PresetType // Optional preset type for specialized setup
	server       *http.Server
	serverCancel context.CancelFunc
	serverWg     *sync.WaitGroup
	shutdownOnce sync.Once
}

// Name implements genkit.Plugin.
func (ot *OpenTelemetry) Name() string {
	return providerID
}

// GetConfig returns the current configuration (mainly for testing purposes).
func (ot *OpenTelemetry) GetConfig() Config {
	return ot.config
}

// New creates a new OpenTelemetry plugin with the given config.
func New(config Config) *OpenTelemetry {
	config.setDefaults()
	return &OpenTelemetry{
		config:   config,
		serverWg: &sync.WaitGroup{},
	}
}

// Init initializes the OpenTelemetry plugin.
func (ot *OpenTelemetry) Init(ctx context.Context) []api.Action {
	// Check if we should export in dev environment
	shouldExport := ot.config.ForceExport || os.Getenv("GENKIT_ENV") != "dev"
	if !shouldExport {
		return nil
	}

	// Initialize trace exporter
	if err := ot.setupTracing(ctx); err != nil {
		panic(fmt.Sprintf("failed to setup tracing: %v", err))
	}

	// Initialize metric exporter
	if err := ot.setupMetrics(ctx); err != nil {
		panic(fmt.Sprintf("failed to setup metrics: %v", err))
	}

	// Initialize log handler
	if err := ot.setupLogging(); err != nil {
		panic(fmt.Sprintf("failed to setup logging: %v", err))
	}

	// Set up signal handling for graceful shutdown if a server was started
	ot.setupSignalHandler()

	return []api.Action{}
}

// setupTracing configures trace export.
func (ot *OpenTelemetry) setupTracing(ctx context.Context) error {
	var spanExporter trace.SpanExporter
	var err error

	if ot.config.TraceExporter != nil {
		spanExporter = ot.config.TraceExporter
	} else {
		spanExporter, err = ot.createDefaultTraceExporter(ctx)
		if err != nil {
			return err
		}
	}

	spanProcessor := trace.NewBatchSpanProcessor(spanExporter)
	tracing.TracerProvider().RegisterSpanProcessor(spanProcessor)

	return nil
}

// setupMetrics configures metric export.
func (ot *OpenTelemetry) setupMetrics(ctx context.Context) error {
	// Use specialized Prometheus setup if this is a Prometheus preset or if forced
	if (ot.presetType != nil && *ot.presetType == PresetPrometheus) || ot.config.EnablePrometheusExporter {
		return ot.setupPrometheusMetrics(ctx)
	}

	var metricExporter metric.Exporter
	var err error

	if ot.config.MetricExporter != nil {
		metricExporter = ot.config.MetricExporter
	} else {
		metricExporter, err = ot.createDefaultMetricExporter(ctx)
		if err != nil {
			return err
		}
	}

	reader := metric.NewPeriodicReader(
		metricExporter,
		metric.WithInterval(ot.config.MetricInterval),
	)

	meterProvider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(meterProvider)

	return nil
}

// setupLogging configures log export.
func (ot *OpenTelemetry) setupLogging() error {
	var handler slog.Handler

	if ot.config.LogHandler != nil {
		handler = ot.config.LogHandler
	} else {
		handler = ot.createDefaultLogHandler()
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

// Shutdown gracefully shuts down the OpenTelemetry plugin and any running servers.
// This method is safe to call multiple times.
func (ot *OpenTelemetry) Shutdown(ctx context.Context) error {
	var shutdownErr error

	ot.shutdownOnce.Do(func() {
		slog.Info("Shutting down OpenTelemetry plugin...")

		// Cancel server context if it exists
		if ot.serverCancel != nil {
			ot.serverCancel()
		}

		// Shutdown HTTP server if it exists
		if ot.server != nil {
			slog.Info("Shutting down Prometheus metrics server...")

			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := ot.server.Shutdown(shutdownCtx); err != nil {
				slog.Error("Error shutting down Prometheus metrics server", "error", err)
				shutdownErr = fmt.Errorf("failed to shutdown Prometheus server: %w", err)
			} else {
				slog.Info("Prometheus metrics server shut down successfully")
			}
		}

		// Wait for server goroutines to finish
		if ot.serverWg != nil {
			done := make(chan struct{})
			go func() {
				ot.serverWg.Wait()
				close(done)
			}()

			select {
			case <-done:
				slog.Info("All OpenTelemetry servers stopped")
			case <-time.After(15 * time.Second):
				slog.Warn("Timeout waiting for OpenTelemetry servers to stop")
			}
		}
	})

	return shutdownErr
}

// setupSignalHandler sets up signal handling for graceful shutdown.
// This should be called after Init to ensure proper cleanup when the application terminates.
func (ot *OpenTelemetry) setupSignalHandler() {
	if ot.server == nil {
		return // No server to manage
	}

	// Create a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start a goroutine to handle signals
	go func() {
		<-sigChan
		slog.Info("Received shutdown signal, starting graceful shutdown...")

		// Create a context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := ot.Shutdown(ctx); err != nil {
			slog.Error("Error during shutdown", "error", err)
		}
	}()
}
