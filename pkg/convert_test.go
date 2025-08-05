package seeker_test

import (
	"github.com/kyma-project/infrastructure-manager/pkg/config"
	"testing"

	gardener_types "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	seeker "github.com/kyma-project/gardener-syncer/pkg"
	"github.com/kyma-project/gardener-syncer/pkg/types"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testProviderType1 = "test-provider-type1"
	testProviderType2 = "test-provider-type2"
	testRegion1       = "test-region1"
	testRegion2       = "test-region2"
	testRegion3       = "test-region3"
	testTaintKey      = "test-key-taint"

	testSeedInDeletion = gardener_types.Seed{
		ObjectMeta: metav1.ObjectMeta{
			DeletionTimestamp: &metav1.Time{},
		},
	}
	testSeedNotVisible = gardener_types.Seed{
		Spec: gardener_types.SeedSpec{
			Settings: &gardener_types.SeedSettings{
				Scheduling: &gardener_types.SeedSettingScheduling{
					Visible: false,
				},
			},
		},
	}
	testSeedNoLatOperation = gardener_types.Seed{
		Spec: gardener_types.SeedSpec{
			Settings: &gardener_types.SeedSettings{
				Scheduling: &gardener_types.SeedSettingScheduling{
					Visible: true,
				},
			},
		},
	}
	testSeedNoSeedGardenletReady = gardener_types.Seed{
		Spec: gardener_types.SeedSpec{
			Settings: &gardener_types.SeedSettings{
				Scheduling: &gardener_types.SeedSettingScheduling{
					Visible: true,
				},
			},
		},
		Status: gardener_types.SeedStatus{
			LastOperation: &gardener_types.LastOperation{},
		},
	}
	testSeedGardenletReadyFalse = gardener_types.Seed{
		Spec: gardener_types.SeedSpec{
			Settings: &gardener_types.SeedSettings{
				Scheduling: &gardener_types.SeedSettingScheduling{
					Visible: true,
				},
			},
		},
		Status: gardener_types.SeedStatus{
			Conditions: []gardener_types.Condition{
				{
					Type:   gardener_types.SeedGardenletReady,
					Status: gardener_types.ConditionFalse,
				},
			},
			LastOperation: &gardener_types.LastOperation{},
		},
	}
	testSeedNoSeedBackupBucketsReady = gardener_types.Seed{
		Spec: gardener_types.SeedSpec{
			Backup: &gardener_types.Backup{},
			Settings: &gardener_types.SeedSettings{
				Scheduling: &gardener_types.SeedSettingScheduling{
					Visible: true,
				},
			},
		},
		Status: gardener_types.SeedStatus{
			Conditions: []gardener_types.Condition{
				{
					Type:   gardener_types.SeedGardenletReady,
					Status: gardener_types.ConditionTrue,
				},
			},
			LastOperation: &gardener_types.LastOperation{},
		},
	}
	testSeedSeedBackupBucketsReadyFalse = gardener_types.Seed{
		Spec: gardener_types.SeedSpec{
			Backup: &gardener_types.Backup{},
			Settings: &gardener_types.SeedSettings{
				Scheduling: &gardener_types.SeedSettingScheduling{
					Visible: true,
				},
			},
		},
		Status: gardener_types.SeedStatus{
			Conditions: []gardener_types.Condition{
				{
					Type:   gardener_types.SeedGardenletReady,
					Status: gardener_types.ConditionTrue,
				},
				{
					Type:   gardener_types.SeedBackupBucketsReady,
					Status: gardener_types.ConditionFalse,
				},
			},
			LastOperation: &gardener_types.LastOperation{},
		},
	}
	testSeedWithTaints = gardener_types.Seed{
		Spec: gardener_types.SeedSpec{
			Provider: gardener_types.SeedProvider{
				Type:   testProviderType1,
				Region: testRegion3,
			},
			Settings: &gardener_types.SeedSettings{
				Scheduling: &gardener_types.SeedSettingScheduling{
					Visible: true,
				},
			},
			Taints: []gardener_types.SeedTaint{
				{
					Key: testTaintKey,
				},
			},
		},
		Status: gardener_types.SeedStatus{
			Conditions: []gardener_types.Condition{
				{
					Type:   gardener_types.SeedGardenletReady,
					Status: gardener_types.ConditionTrue,
				},
			},
			LastOperation: &gardener_types.LastOperation{},
		},
	}

	testSeedWithToleratedTaints = gardener_types.Seed{
		Spec: gardener_types.SeedSpec{
			Provider: gardener_types.SeedProvider{
				Type:   testProviderType1,
				Region: testRegion1,
			},
			Settings: &gardener_types.SeedSettings{
				Scheduling: &gardener_types.SeedSettingScheduling{
					Visible: true,
				},
			},
			Taints: []gardener_types.SeedTaint{
				{
					Key: testTaintKey,
				},
			},
		},
		Status: gardener_types.SeedStatus{
			Conditions: []gardener_types.Condition{
				{
					Type:   gardener_types.SeedGardenletReady,
					Status: gardener_types.ConditionTrue,
				},
			},
			LastOperation: &gardener_types.LastOperation{},
		},
	}

	testSeedOK = gardener_types.Seed{
		Spec: gardener_types.SeedSpec{
			Provider: gardener_types.SeedProvider{
				Type:   testProviderType1,
				Region: testRegion1,
			},
			Settings: &gardener_types.SeedSettings{
				Scheduling: &gardener_types.SeedSettingScheduling{
					Visible: true,
				},
			},
		},
		Status: gardener_types.SeedStatus{
			Conditions: []gardener_types.Condition{
				{
					Type:   gardener_types.SeedGardenletReady,
					Status: gardener_types.ConditionTrue,
				},
			},
			LastOperation: &gardener_types.LastOperation{},
		},
	}
	testSeedOKWithBackup = gardener_types.Seed{
		Spec: gardener_types.SeedSpec{
			Provider: gardener_types.SeedProvider{
				Type:   testProviderType2,
				Region: testRegion2,
			},
			Backup: &gardener_types.Backup{},
			Settings: &gardener_types.SeedSettings{
				Scheduling: &gardener_types.SeedSettingScheduling{
					Visible: true,
				},
			},
		},
		Status: gardener_types.SeedStatus{
			Conditions: []gardener_types.Condition{
				{
					Type:   gardener_types.SeedGardenletReady,
					Status: gardener_types.ConditionTrue,
				},
				{
					Type:   gardener_types.SeedBackupBucketsReady,
					Status: gardener_types.ConditionTrue,
				},
			},
			LastOperation: &gardener_types.LastOperation{},
		},
	}
)

func TestToProvideRegions(t *testing.T) {

	testCases := []struct {
		name     string
		seeds    []gardener_types.Seed
		expected types.Providers
	}{
		{
			name: "seed filtered",
			seeds: []gardener_types.Seed{
				testSeedInDeletion,
				testSeedNotVisible,
				testSeedNoLatOperation,
				testSeedNoSeedGardenletReady,
				testSeedGardenletReadyFalse,
				testSeedNoSeedBackupBucketsReady,
				testSeedSeedBackupBucketsReadyFalse,
				testSeedWithTaints,
			},
			expected: types.Providers{},
		},
		{
			name: "tolerations",
			seeds: []gardener_types.Seed{
				testSeedWithToleratedTaints,
				testSeedWithTaints,
			},
			expected: types.Providers{
				testSeedWithToleratedTaints.Spec.Provider.Type: {
					SeedRegions: []string{
						testSeedWithToleratedTaints.Spec.Provider.Region,
					},
				},
			},
		},
		{
			name: "seed found",
			seeds: []gardener_types.Seed{
				testSeedOKWithBackup,
			},
			expected: types.Providers{
				testSeedOKWithBackup.Spec.Provider.Type: {
					SeedRegions: []string{
						testSeedOKWithBackup.Spec.Provider.Region,
					},
				},
			},
		},
		{
			name: "mixed",
			seeds: []gardener_types.Seed{
				testSeedInDeletion,
				testSeedNotVisible,
				testSeedNoLatOperation,
				testSeedNoSeedGardenletReady,
				testSeedGardenletReadyFalse,
				testSeedNoSeedBackupBucketsReady,
				testSeedSeedBackupBucketsReadyFalse,
				testSeedWithTaints,
				testSeedOKWithBackup,
				testSeedOK,
			},
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
			// WHEN
			actual := seeker.ToProviderRegions(testCase.seeds, config.TolerationsConfig{
				testRegion1: {{Key: testTaintKey}},
			})

			// THEN
			require.Equal(t, testCase.expected, actual)
		})
	}
}
