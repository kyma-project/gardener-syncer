package types_test

import (
	"testing"

	"github.com/kyma-project/gardener-syncer/pkg/types"
	"github.com/stretchr/testify/require"
)

var (
	testProviderName = "test-provider"
	testRegionName   = "test-region"
	testSeed         = "test-seed"
)

func TestProviderRegionSeed_Add(t *testing.T) {
	type data struct {
		provider string
		region   string
	}

	testCases := []struct {
		name    string
		initial types.Providers
		data    data
	}{
		{
			name:    "empty",
			initial: types.Providers{},
			data: data{
				provider: testProviderName,
				region:   testRegionName,
			},
		},
		{
			name: "initialised",
			initial: types.Providers{
				testProviderName: {
					SeedRegions: []string{"some-other-test-seed"},
				},
			},
			data: data{
				provider: testProviderName,
				region:   testRegionName,
			},
		},
		{
			name: "duplicated",
			initial: types.Providers{
				testProviderName: {
					SeedRegions: []string{testRegionName},
				},
			},
			data: data{
				provider: testProviderName,
				region:   testRegionName,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// WHEN
			testCase.initial.Add(
				testCase.data.provider,
				testCase.data.region,
			)
			// THEN
			regions := testCase.initial[testCase.data.provider].SeedRegions
			require.Contains(t, regions, testCase.data.region)
		})
	}
}
