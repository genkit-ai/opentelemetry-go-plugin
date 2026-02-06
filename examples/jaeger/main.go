package main

import (
	"context"
	"log"
	"time"

	"github.com/firebase/genkit/go/genkit"
	opentelemetry "github.com/xavidop/genkit-opentelemetry-go"
)
// test
func main() {
	ctx := context.Background()

	jaegerExample(ctx)

	time.Sleep(120 * time.Minute) // Allow time for exporters to flush
}

func jaegerExample(ctx context.Context) {

	// Example: Using Jaeger preset (commented out to avoid interference)
	jaegerPlugin := opentelemetry.NewWithPreset(opentelemetry.PresetJaeger, opentelemetry.Config{
		ServiceName: "my-genkit-app",
		ForceExport: true, // Export even in development
	})

	// Initialize Genkit
	genkit.Init(ctx,
		genkit.WithPlugins(jaegerPlugin),
	)

	log.Println("Preset examples completed")
}
