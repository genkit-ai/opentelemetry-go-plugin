# Genkit OpenTelemetry Plugin for Go

A comprehensive OpenTelemetry plugin for [Genkit Go](https://genkit.dev/go/docs/) that provides out-of-the-box trace spans, metrics collectors, and logs with the flexibility to bring your own exporters.

## Features

- 🚀 **Out-of-the-box setup** with sensible defaults
- 📊 **Multiple exporters** - OTLP, Prometheus, Jaeger, Console, and more
- ⚙️ **Flexible configuration** - use presets or bring your own exporters
- 🔄 **Environment-aware** - different behavior for dev vs production
- 📈 **Comprehensive telemetry** - traces, metrics, and structured logs

## Installation

```bash
go get github.com/xavidop/genkit-opentelemetry-go
```

## Quick Start

### Basic Usage

The simplest way to get started is with the default OTLP configuration:

```go
package main

import (
    "context"
    "log"

    "github.com/firebase/genkit/go/genkit"
    opentelemetry "github.com/xavidop/genkit-opentelemetry-go"
)

func main() {
    ctx := context.Background()

    // Initialize OpenTelemetry plugin with default settings
    plugin := opentelemetry.New(opentelemetry.Config{
        ServiceName: "my-genkit-app",
        ForceExport: true, // Export even in development
    })

    // Initialize Genkit
    genkit.Init(ctx,
        genkit.WithPlugins(plugin),
    )


    // Your Genkit flows will now automatically emit telemetry!
}
```

### Development Setup

For development, use the console preset to see telemetry data in your terminal:

```go
otelPlugin := opentelemetry.NewWithPreset(opentelemetry.PresetConsole, opentelemetry.Config{
    ServiceName: "my-app-dev",
    LogLevel:    slog.LevelDebug,
    ForceExport: true, // Export even in dev environment
})
```

## Examples

### Using Jaeger

```go
jaegerPlugin := opentelemetry.NewWithPreset(opentelemetry.PresetJaeger, opentelemetry.Config{
    ServiceName: "my-genkit-app",
    ForceExport: true,
})
```

### Using Prometheus

```go
plugin := opentelemetry.NewWithPreset(opentelemetry.PresetPrometheus, opentelemetry.Config{
    ServiceName:    "my-genkit-app",
    ForceExport:    true,
    MetricInterval: 15 * time.Second,
    PrometheusPort: 8081, // Custom port for the /metrics endpoint
})
```

### Using OTLP

```go
plugin := opentelemetry.NewWithPreset(opentelemetry.PresetOTLP, opentelemetry.Config{
    ServiceName:    "my-genkit-app",
    ForceExport:    true,
    MetricInterval: 15 * time.Second,
})
```

## Production Setup Options

### Generic OTLP Backend

```go
otelPlugin := opentelemetry.New(opentelemetry.Config{
    ServiceName:  "my-app",
    OTLPEndpoint: "https://api.your-provider.com",
    OTLPUseHTTP:  true,
    OTLPHeaders: map[string]string{
        "authorization": "Bearer your-token",
    },
})
```

### Popular Observability Providers

#### Honeycomb

```go
otelPlugin := opentelemetry.New(opentelemetry.Config{
    ServiceName:  "my-app",
    OTLPEndpoint: "https://api.honeycomb.io",
    OTLPUseHTTP:  true,
    OTLPHeaders: map[string]string{
        "x-honeycomb-team": os.Getenv("HONEYCOMB_API_KEY"),
    },
})
```

#### Datadog

```go
otelPlugin := opentelemetry.New(opentelemetry.Config{
    ServiceName:  "my-app",
    OTLPEndpoint: "https://trace.agent.datadoghq.com",
    OTLPUseHTTP:  true,
    OTLPHeaders: map[string]string{
        "dd-api-key": os.Getenv("DD_API_KEY"),
    },
})
```

#### New Relic

```go
otelPlugin := opentelemetry.New(opentelemetry.Config{
    ServiceName:  "my-app",
    OTLPEndpoint: "https://otlp.nr-data.net",
    OTLPUseHTTP:  true,
    OTLPHeaders: map[string]string{
        "api-key": os.Getenv("NEW_RELIC_LICENSE_KEY"),
    },
})
```

#### Langfuse

```go
// Build basic auth header value: base64("public_key:secret_key")
auth := base64.StdEncoding.EncodeToString(
	[]byte(os.Getenv("LANGFUSE_PUBLIC_KEY") + ":" + os.Getenv("LANGFUSE_SECRET_KEY")),
)

otelPlugin := opentelemetry.New(opentelemetry.Config{
    ServiceName:  "my-app",
    OTLPEndpoint: "https://cloud.langfuse.com/api/public/otel/v1/traces", // EU region
    // OTLPEndpoint: "https://us.cloud.langfuse.com/api/public/otel/v1/traces", // US region
    OTLPUseHTTP:  true,
    OTLPHeaders: map[string]string{
        "Authorization": "Basic " + auth,
    },
})
```

#### Laminar

```go
otelPlugin := opentelemetry.New(opentelemetry.Config{
    DisableMetricsExporter: true,
    ServiceName:  "my-app",
    OTLPEndpoint: "https://api.lmnr.ai:8443",
    OTLPHeaders: map[string]string{
        "Authorization": "Bearer " + os.Getenv("LAMINAR_API_KEY"),
    },
})
```

## Local Testing with Docker

### 1. Start the Observability Stack

```bash
# Clone the repository
git clone https://github.com/xavidop/genkit-opentelemetry-go
cd genkit-opentelemetry-go

# Start Jaeger, Prometheus, and Grafana
docker-compose up -d
```

### 2. Run Your Application

```bash
# Use the OTLP collector with gRPC (default)
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 go run your-app.go

# Use the OTLP collector with HTTP
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318 go run your-app.go

# Or run the provided examples
go run examples/basic/main.go
go run examples/jaeger/main.go  
go run examples/prometheus/main.go
go run examples/otel/main.go
```

### 3. View Your Data

- **Jaeger UI**: http://localhost:16686 (traces)
- **Prometheus**: http://localhost:9090 (metrics)
- **Grafana**: http://localhost:3000 (dashboards, admin/admin)

## Environment Variables

The plugin respects standard OpenTelemetry environment variables:

```bash
# OTLP Endpoint - for gRPC (default), use host:port format
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317

# OTLP Endpoint - for HTTP, use full URL format  
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318

# Service name
export OTEL_SERVICE_NAME=my-genkit-app

# Headers
export OTEL_EXPORTER_OTLP_HEADERS=authorization=Bearer-token

# Protocol (grpc or http/protobuf)
export OTEL_EXPORTER_OTLP_PROTOCOL=grpc

# Special endpoints for stdout
export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=stdout
export OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=stdout
```


## Configuration Options

### Config Structure

```go
type Config struct {
    // Export even in the dev environment
    ForceExport bool

    // The interval for exporting metric data (default: 60s)
    MetricInterval time.Duration

    // The minimum level at which logs will be written (default: Info)
    LogLevel slog.Leveler

    // Custom trace span exporter (optional)
    TraceExporter trace.SpanExporter

    // Custom metric exporter (optional)
    MetricExporter metric.Exporter

    // Custom log handler (optional)
    LogHandler slog.Handler

    // OTLP endpoint (default: "localhost:4317")
    OTLPEndpoint string

    // Whether to use HTTP instead of gRPC for OTLP (default: false)
    OTLPUseHTTP bool

    // Headers to include in OTLP requests
    OTLPHeaders map[string]string

    // Service name for telemetry data (default: "genkit-service")
    ServiceName string

    // Service version for telemetry data
    ServiceVersion string

    // Additional resource attributes
    ResourceAttributes map[string]string
}
```

## Presets

The plugin comes with several presets for common scenarios:
### OTLP (Default)
```go
otelPlugin := opentelemetry.NewWithPreset(opentelemetry.PresetOTLP)
```

### Jaeger
```go
otelPlugin := opentelemetry.NewWithPreset(opentelemetry.PresetJaeger, opentelemetry.Config{
    OTLPEndpoint: "http://localhost:14268/api/traces",
})
```

### Prometheus
```go
otelPlugin := opentelemetry.NewWithPreset(opentelemetry.PresetPrometheus)
```

### Console (Development)
```go
otelPlugin := opentelemetry.NewWithPreset(opentelemetry.PresetConsole, opentelemetry.Config{
    LogLevel: slog.LevelDebug,
})
```

## What's Automatically Instrumented

When you use this plugin with Genkit, you automatically get:

- **Traces** for all Genkit flows and actions
- **Metrics** for flow execution times, success/failure rates
- **Logs** with proper correlation to traces
- **Custom attributes** for Genkit-specific metadata

## Adding Custom Instrumentation

You can add your own traces and metrics:

```go
import "go.opentelemetry.io/otel"

// Get a tracer
tracer := otel.Tracer("my-component")

// Create a span
ctx, span := tracer.Start(ctx, "my-operation")
defer span.End()

// Add attributes
span.SetAttributes(
    attribute.String("user.id", userID),
    attribute.Int("batch.size", batchSize),
)

// Get a meter for custom metrics
meter := otel.Meter("my-component")
counter, _ := meter.Int64Counter("requests_total")
counter.Add(ctx, 1, metric.WithAttributes(
    attribute.String("method", "POST"),
))
```

## Troubleshooting

### Problem: No traces appearing

1. Check if `ForceExport: true` is set in development
2. Verify the OTLP endpoint is correct
3. Check network connectivity to your telemetry backend
4. Look for error logs in your application

### Problem: High memory usage

1. Reduce the `MetricInterval` for more frequent exports
2. Consider using sampling for traces in high-traffic applications
3. Check if your telemetry backend is accepting data

### Problem: Missing data

1. Check if your telemetry backend has ingestion limits
2. Verify the service name matches your expectations

### Getting Help

1. Check the [examples](./examples/) directory
2. Review the [main documentation](./README.md)
3. Look at the test files for usage patterns
4. Open an issue if you find bugs or need features

## Advanced Usage

### Custom Exporters

You can bring your own exporters:

```go
import (
    "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
)

// Create custom exporters
traceExporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
metricExporter, _ := stdoutmetric.New(stdoutmetric.WithPrettyPrint())

otelPlugin := opentelemetry.New(opentelemetry.Config{
    ServiceName:     "my-app",
    TraceExporter:   traceExporter,
    MetricExporter:  metricExporter,
    ForceExport:     true,
})
```

## Contributing

Check out the [CONTRIBUTING.md](CONTRIBUTING.md) file for details on how to contribute to this project.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [Genkit](https://genkit.dev/) - AI application framework
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/) - Observability framework
- [Genkit Google Cloud Plugin](https://genkit.dev/go/docs/plugins/google-cloud/) - Official Google Cloud telemetry plugin
