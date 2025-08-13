package seeker_test

import (
	"context"
	"fmt"
	"github.com/kyma-project/infrastructure-manager/pkg/config"
	"testing"

	gardener_types "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	seeker "github.com/kyma-project/gardener-syncer/pkg"
	"github.com/kyma-project/gardener-syncer/pkg/types"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	errListFailedTest = fmt.Errorf("list failed test")
)

func TestBuildFet(t *testing.T) {
	testCases := []struct {
		name        string
		expected    types.Providers
		list        seeker.List
		expectedErr error
	}{
		{
			name:        "list error",
			list:        buildListWithError(errFetchSeedsFailedTest),
			expectedErr: errFetchSeedsFailedTest,
		},
		{
			name:     "list empty",
			list:     buildList(gardener_types.SeedList{}),
			expected: types.Providers{},
		},
		{
			name: "OK",
			list: buildList(gardener_types.SeedList{
				TypeMeta: metav1.TypeMeta{},
				Items: []gardener_types.Seed{
					testSeedOK,
					testSeedOKWithBackup,
				},
			}),
			expected: types.Providers{
				testSeedOK.Spec.Provider.Type: {
					SeedRegions: []string{
						testSeedOK.Spec.Provider.Region,
					},
				},
				testSeedOKWithBackup.Spec.Provider.Type: {
					SeedRegions: []string{
						testSeedOKWithBackup.Spec.Provider.Region,
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// GIVEN
			fetchSeeds := seeker.BuildFetchSeedFn(seeker.FetchSeedsOpts{
				List: testCase.list,
			}, config.TolerationsConfig{})

			// WHEN
			actual, err := fetchSeeds()

			// THEN
			if testCase.expectedErr != nil {
				require.EqualError(t, err, testCase.expectedErr.Error())
				return
			}

			// THEN
			require.NoError(t, err)
			require.Equal(t, testCase.expected, actual)
		})
	}
}

func buildListWithError(err error) seeker.List {
	return func(context.Context, client.ObjectList, ...client.ListOption) error {
		return err
	}
}

func buildList(out gardener_types.SeedList) seeker.List {
	return func(ctx context.Context, ol client.ObjectList, lo ...client.ListOption) error {
		seedList := ol.(*gardener_types.SeedList)
		*seedList = out
		return nil
	}
}
