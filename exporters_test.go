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
	"testing"
)

func TestEndpointParse(t *testing.T) {
	tests := []struct {
		name         string
		endpoint     string
		wantScheme   string
		wantHostPort string
		wantPath     string
		wantUseTLS   bool
	}{
		{
			name:         "host:port only",
			endpoint:     "localhost:4317",
			wantScheme:   "",
			wantHostPort: "localhost:4317",
			wantPath:     "",
			wantUseTLS:   false,
		},
		{
			name:         "http scheme without path",
			endpoint:     "http://localhost:4318",
			wantScheme:   "http",
			wantHostPort: "localhost:4318",
			wantPath:     "",
			wantUseTLS:   false,
		},
		{
			name:         "https scheme without path",
			endpoint:     "https://otel.example.com:443",
			wantScheme:   "https",
			wantHostPort: "otel.example.com:443",
			wantPath:     "",
			wantUseTLS:   true,
		},
		{
			name:         "https with custom path",
			endpoint:     "https://tracing.example.com:443/otel/v1/traces",
			wantScheme:   "https",
			wantHostPort: "tracing.example.com:443",
			wantPath:     "/otel/v1/traces",
			wantUseTLS:   true,
		},
		{
			name:         "host:port with path",
			endpoint:     "localhost:4318/custom/path",
			wantScheme:   "",
			wantHostPort: "localhost:4318/custom/path", // No scheme, treat entire input as HostPort
			wantPath:     "",
			wantUseTLS:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var endpoint EndpointURL
			endpoint.Parse(tt.endpoint)

			if endpoint.Scheme != tt.wantScheme {
				t.Errorf("Scheme = %q, want %q", endpoint.Scheme, tt.wantScheme)
			}
			if endpoint.HostPort != tt.wantHostPort {
				t.Errorf("HostPort = %q, want %q", endpoint.HostPort, tt.wantHostPort)
			}
			if endpoint.Path != tt.wantPath {
				t.Errorf("Path = %q, want %q", endpoint.Path, tt.wantPath)
			}
			if endpoint.UseTLS != tt.wantUseTLS {
				t.Errorf("UseTLS = %v, want %v", endpoint.UseTLS, tt.wantUseTLS)
			}
		})
	}
}
