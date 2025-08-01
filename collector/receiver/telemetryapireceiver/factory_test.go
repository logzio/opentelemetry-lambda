// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetryapireceiver // import "github.com/open-telemetry/opentelemetry-lambda/collector/receiver/telemetryapireceiver"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestNewFactory(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "creates a new factory with correct type",
			testFunc: func(t *testing.T) {
				factory := NewFactory("test")
				require.EqualValues(t, typeStr, factory.Type().String())
			},
		},
		{
			desc: "creates a new factory with valid default config",
			testFunc: func(t *testing.T) {
				factory := NewFactory("test")

				var expectedCfg component.Config = &Config{
					extensionID: "test",
					Port:        defaultPort,
					Types:       []string{platform, function, extension},
					MaxItems:    defaultMaxItems,
					MaxBytes:    defaultMaxBytes,
					TimeoutMS:   defaultTimeoutMS,
				}

				require.Equal(t, expectedCfg, factory.CreateDefaultConfig())
			},
		},
		{
			desc: "creates a new factory and CreateTracesReceiver returns no error",
			testFunc: func(t *testing.T) {
				factory := NewFactory("test")
				cfg := factory.CreateDefaultConfig()
				_, err := factory.CreateTraces(
					context.Background(),
					receivertest.NewNopSettings(Type),
					cfg,
					consumertest.NewNop(),
				)
				require.NoError(t, err)
			},
		},
		{
			desc: "creates a new factory and CreateTracesReceiver returns error with incorrect config",
			testFunc: func(t *testing.T) {
				factory := NewFactory("test")
				_, err := factory.CreateTraces(
					context.Background(),
					receivertest.NewNopSettings(Type),
					nil,
					consumertest.NewNop(),
				)
				require.ErrorIs(t, err, errConfigNotTelemetryAPI)
			},
		},
		{
			desc: "creates a new factory and CreateLogsReceiver returns no error",
			testFunc: func(t *testing.T) {
				factory := NewFactory("test")
				cfg := factory.CreateDefaultConfig()
				_, err := factory.CreateLogs(
					context.Background(),
					receivertest.NewNopSettings(Type),
					cfg,
					consumertest.NewNop(),
				)
				require.NoError(t, err)
			},
		},
		{
			desc: "creates a new factory and CreateMetricsReceiver returns no error",
			testFunc: func(t *testing.T) {
				factory := NewFactory("test")
				cfg := factory.CreateDefaultConfig()
				_, err := factory.CreateMetrics(
					context.Background(),
					receivertest.NewNopSettings(Type),
					cfg,
					consumertest.NewNop(),
				)
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}
