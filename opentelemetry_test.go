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
	"sync"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestShouldSetupExporters(t *testing.T) {
	testCases := []struct {
		name                   string
		config                 Config
		wantShouldSetupTracing bool
		wantShouldSetupMetrics bool
	}{
		{
			name:                   "default enabled",
			config:                 Config{},
			wantShouldSetupTracing: true,
			wantShouldSetupMetrics: true,
		},
		{
			name: "all disabled",
			config: Config{
				DisableTracingExporter: true,
				DisableMetricsExporter: true,
			},
			wantShouldSetupTracing: false,
			wantShouldSetupMetrics: false,
		},
		{
			name: "tracing disabled only",
			config: Config{
				DisableTracingExporter: true,
			},
			wantShouldSetupTracing: false,
			wantShouldSetupMetrics: true,
		},
		{
			name: "metrics disabled only",
			config: Config{
				DisableMetricsExporter: true,
			},
			wantShouldSetupTracing: true,
			wantShouldSetupMetrics: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ot := &OpenTelemetry{config: tc.config}
			if got := ot.shouldSetupTracing(); got != tc.wantShouldSetupTracing {
				t.Errorf("shouldSetupTracing() = %v, want %v", got, tc.wantShouldSetupTracing)
			}
			if got := ot.shouldSetupMetrics(); got != tc.wantShouldSetupMetrics {
				t.Errorf("shouldSetupMetrics() = %v, want %v", got, tc.wantShouldSetupMetrics)
			}
		})
	}
}

type captureSpanExporter struct {
	mu    sync.Mutex
	spans []trace.ReadOnlySpan
}

func (e *captureSpanExporter) ExportSpans(_ context.Context, spans []trace.ReadOnlySpan) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.spans = append(e.spans, spans...)
	return nil
}

func (e *captureSpanExporter) Shutdown(context.Context) error {
	return nil
}

func (e *captureSpanExporter) spansSnapshot() []trace.ReadOnlySpan {
	e.mu.Lock()
	defer e.mu.Unlock()

	out := make([]trace.ReadOnlySpan, len(e.spans))
	copy(out, e.spans)
	return out
}

func TestSetupTracingSetsConfiguredResourceAttributesOnExportedSpans(t *testing.T) {
	t.Cleanup(func() {
		otel.SetTracerProvider(trace.NewTracerProvider())
	})

	exporter := &captureSpanExporter{}
	ot := New(Config{
		ForceExport:    true,
		ServiceName:    "checkout-service",
		ServiceVersion: "1.2.3",
		ResourceAttributes: map[string]string{
			"deployment.environment": "test",
		},
		TraceExporter: exporter,
	})

	res, err := ot.createResource(context.Background())
	if err != nil {
		t.Fatalf("createResource() error = %v", err)
	}

	if err := ot.setupTracing(context.Background(), res); err != nil {
		t.Fatalf("setupTracing() error = %v", err)
	}

	_, span := otel.Tracer("test").Start(context.Background(), "test-span")
	span.End()

	tp, ok := otel.GetTracerProvider().(*trace.TracerProvider)
	if !ok {
		t.Fatalf("global tracer provider is %T, want *trace.TracerProvider", otel.GetTracerProvider())
	}
	if err := tp.ForceFlush(context.Background()); err != nil {
		t.Fatalf("ForceFlush() error = %v", err)
	}

	spans := exporter.spansSnapshot()
	if len(spans) == 0 {
		t.Fatal("expected at least one exported span")
	}

	t.Run("service name", func(t *testing.T) {
		serviceName, ok := resourceAttributeValue(spans[0], "service.name")
		if !ok {
			t.Fatal("service.name attribute not found on span resource")
		}
		if serviceName != "checkout-service" {
			t.Fatalf("service.name = %q, want %q", serviceName, "checkout-service")
		}
	})

	t.Run("service version", func(t *testing.T) {
		serviceVersion, ok := resourceAttributeValue(spans[0], "service.version")
		if !ok {
			t.Fatal("service.version attribute not found on span resource")
		}
		if serviceVersion != "1.2.3" {
			t.Fatalf("service.version = %q, want %q", serviceVersion, "1.2.3")
		}
	})

	t.Run("resource attributes", func(t *testing.T) {
		environment, ok := resourceAttributeValue(spans[0], "deployment.environment")
		if !ok {
			t.Fatal("deployment.environment attribute not found on span resource")
		}
		if environment != "test" {
			t.Fatalf("deployment.environment = %q, want %q", environment, "test")
		}
	})
}

func resourceAttributeValue(span trace.ReadOnlySpan, key string) (string, bool) {
	for _, attr := range span.Resource().Attributes() {
		if string(attr.Key) == key {
			return attr.Value.AsString(), true
		}
	}
	return "", false
}
