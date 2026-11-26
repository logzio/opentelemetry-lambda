// Copyright The OpenTelemetry Authors
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

package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/confmap/xconfmap"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/open-telemetry/opentelemetry-lambda/collector/lambdacomponents"
)

// The shipped configs must load with this collector's components: an unknown component name
// (for example a receiver registered under a different type) fails here instead of at runtime.
func TestShippedConfigsAreValid(t *testing.T) {
	t.Setenv("LOGZIO_GATEWAY_ENDPOINT", "https://gateway.example.com")
	t.Setenv("LOGZIO_GATEWAY_TOKEN", "test-token")

	factories, err := lambdacomponents.Components("collector")
	require.NoError(t, err)

	for _, file := range []string{"config.gateway.yaml"} {
		t.Run(file, func(t *testing.T) {
			provider, err := otelcol.NewConfigProvider(otelcol.ConfigProviderSettings{
				ResolverSettings: confmap.ResolverSettings{
					URIs:              []string{"file:" + file},
					ProviderFactories: []confmap.ProviderFactory{fileprovider.NewFactory(), envprovider.NewFactory()},
				},
			})
			require.NoError(t, err)

			cfg, err := provider.Get(context.Background(), factories)
			require.NoError(t, err)
			require.NoError(t, xconfmap.Validate(cfg))
		})
	}
}
