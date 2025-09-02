package seeker

import (
	"log/slog"
	"strings"
	"time"

	"github.com/kyma-project/infrastructure-manager/pkg/config"

	"sigs.k8s.io/yaml"

	gardener_types "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	"github.com/kyma-project/gardener-syncer/pkg/types"
)

func VerifySeedReadiness(seed *gardener_types.Seed) bool {
	if seed.Status.LastOperation == nil {
		return false
	}

	if cond := v1beta1helper.GetCondition(seed.Status.Conditions, gardener_types.SeedGardenletReady); cond == nil || cond.Status != gardener_types.ConditionTrue {
		return false
	}

	if seed.Spec.Backup != nil {
		if cond := v1beta1helper.GetCondition(seed.Status.Conditions, gardener_types.SeedBackupBucketsReady); cond == nil || cond.Status != gardener_types.ConditionTrue {
			return false
		}
	}

	return true
}

func VerifySeedTaints(seed *gardener_types.Seed, tolerationConfig config.TolerationsConfig) bool {
	if len(seed.Spec.Taints) == 0 {
		return true
	}

	tolerations, seedRegionHasTolerations := tolerationConfig[seed.Spec.Provider.Region]

	if !seedRegionHasTolerations {
		return false // If seed has taints and there are no tolerations for the seed region, we cannot use the seed
	}

	for _, taint := range seed.Spec.Taints {
		matched := TaintMatched(taint, tolerations)
		if !matched {
			return false // If any taint does not match its toleration, we cannot use the seed
		}
	}
	return true
}

func TaintMatched(taint gardener_types.SeedTaint, tolerations []gardener_types.Toleration) bool {
	for _, toleration := range tolerations {
		if taint.Key != toleration.Key {
			continue
		}

		if toleration.Value == nil && taint.Value == nil {
			return true // value `nil` only matches `nil` (?)
		}

		if toleration.Value != nil && taint.Value != nil && *taint.Value == *toleration.Value {
			return true
		}
	}
	return false
}

func SeedCanBeUsed(seed *gardener_types.Seed, tolerations config.TolerationsConfig) bool {
	hasNoDeletionTimestamp := seed.DeletionTimestamp == nil
	isReady := VerifySeedReadiness(seed)
	isVisible := seed.Spec.Settings != nil &&
		seed.Spec.Settings.Scheduling != nil &&
		seed.Spec.Settings.Scheduling.Visible

	hasCorrectTaintsConfig := VerifySeedTaints(seed, tolerations)

	result := hasNoDeletionTimestamp && isVisible && isReady && hasCorrectTaintsConfig
	if !result {
		slog.Info("seed rejected",
			"name", seed.Name,
			"hasNoDeletionTimestamp", hasNoDeletionTimestamp,
			"isVisible", isVisible,
			"hasCorrectTaintsConfig", hasCorrectTaintsConfig,
			"isReady", isReady)
	}
	return result
}

func ToProviderRegions(seeds []gardener_types.Seed, tolerations config.TolerationsConfig) (out types.Providers) {
	defer LogWithDuration(time.Now(), "conversion complete")

	out = types.Providers{}
	for _, seed := range seeds {
		if SeedCanBeUsed(&seed, tolerations) {
			out.Add(
				seed.Spec.Provider.Type,
				seed.Spec.Provider.Region,
			)
			continue
		}
	}

	return out
}

func ToConfigMap(providerRegions types.Providers) (map[string]string, error) {
	result := map[string]string{}
	for k, v := range providerRegions {
		data, err := yaml.Marshal(v)
		if err != nil {
			return nil, err
		}
		result[k] = strings.TrimRight(string(data), "\n")
	}
	return result, nil
}

type Convert[T any, V any] func(T) (V, error)
