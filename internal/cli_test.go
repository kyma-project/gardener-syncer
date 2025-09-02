package cli

import (
	"fmt"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	seeker "github.com/kyma-project/gardener-syncer/pkg"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"sigs.k8s.io/yaml"
	"testing"
)

const seedsFilePath = "config/test/seeds_minimal.yaml"
const converterConfigPath = "config/test/converter_config.yaml"

func TestMarshalingStubData(t *testing.T) {
	t.Run("proper marshaling of infrastructure manager config", func(t *testing.T) {
		converter_config, err := loadConverterConfig(converterConfigPath)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if len(converter_config.ConverterConfig.Tolerations) == 0 {
			t.Fatalf("no tolerations found in config")
		}
	})

	t.Run("proper marshaling of seeds example", func(t *testing.T) {
		seeds, err := loadSeeds(seedsFilePath)
		if err != nil {
			t.Fatalf("failed to load seed file: %v", err)
		}

		if seeds.Items == nil {
			t.Fatalf("failed to unmarshal seeds from file: %v", err)
		}
	})
}

func loadSeeds(path string) (seeds v1beta1.SeedList, err error) {
	seedsExamplesFile, err := os.Open(path)
	if err != nil {
		return seeds, fmt.Errorf("unable to open seeds stub data %s: %w", path, err)
	}
	defer seedsExamplesFile.Close()

	yamlData, err := io.ReadAll(seedsExamplesFile)
	yaml.Unmarshal(yamlData, &seeds)

	if err != nil {
		return seeds, fmt.Errorf("unable to open seeds stub data %s: %w", path, err)
	}

	return
}

func TestConfigurationIntegration(t *testing.T) {
	testCases := []struct {
		seedName              string
		expectedReadiness     bool
		expectedSeedTaints    bool
		expectedTaintMatched  bool
		expectedSeedCanBeUsed bool
	}{
		{
			seedName:              "aws-ap1",
			expectedReadiness:     true,
			expectedSeedTaints:    false,
			expectedTaintMatched:  false,
			expectedSeedCanBeUsed: false,
		},
		{
			seedName:              "aws-ap2",
			expectedReadiness:     true,
			expectedSeedTaints:    true,
			expectedTaintMatched:  true,
			expectedSeedCanBeUsed: false,
		},
		{
			seedName:              "gcp-ha-sa1",
			expectedReadiness:     true,
			expectedSeedTaints:    true,
			expectedTaintMatched:  true,
			expectedSeedCanBeUsed: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.seedName, func(t *testing.T) {
			converter_config, _ := loadConverterConfig(converterConfigPath)
			seeds, _ := loadSeeds(seedsFilePath)

			tolerations := converter_config.ConverterConfig.Tolerations

			regionsWithTolerations := make([]string, 0, len(tolerations))

			for s := range tolerations {
				regionsWithTolerations = append(regionsWithTolerations, s)
			}

			for _, region := range regionsWithTolerations {
				for _, item := range seeds.Items {
					if item.Name != testCase.seedName {
						continue
					}

					assert.Equal(t, testCase.expectedSeedTaints, seeker.VerifySeedTaints(&item, tolerations))
					if item.Spec.Taints != nil {
						assert.Equal(t, testCase.expectedTaintMatched, seeker.TaintMatched(item.Spec.Taints[0], tolerations[region]))
					}
					assert.Equal(t, testCase.expectedReadiness, seeker.VerifySeedReadiness(&item))
					assert.Equal(t, testCase.expectedSeedCanBeUsed, seeker.SeedCanBeUsed(&item, tolerations))
				}
			}
		})
	}
}
