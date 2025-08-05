package seeker

import (
	"github.com/kyma-project/infrastructure-manager/pkg/config"
	"log/slog"
	"strings"
	"time"

	"sigs.k8s.io/yaml"

	gardener_types "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	"github.com/kyma-project/gardener-syncer/pkg/types"
)

func verifySeedReadiness(seed *gardener_types.Seed) bool {
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

func checkTolerations(seed *gardener_types.Seed, tolerationConfig config.TolerationsConfig) bool {
	tolerations, regionHasTolerations := tolerationConfig[seed.Spec.Provider.Region]
	taintsN := len(seed.Spec.Taints)
	if !regionHasTolerations || taintsN == 0 {
		return taintsN == 0
	}

tnt:
	for _, taint := range seed.Spec.Taints {
		for _, toleration := range tolerations {
			if taint.Key == toleration.Key && (toleration.Value == nil || taint.Value == toleration.Value) {
				continue tnt
			}
		}
		return false
	}
	return true
}

func seedCanBeUsed(seed *gardener_types.Seed, tolerations config.TolerationsConfig) bool {
	isDeletionTimesampt := seed.DeletionTimestamp == nil
	isReady := verifySeedReadiness(seed)
	isVisible := seed.Spec.Settings != nil &&
		seed.Spec.Settings.Scheduling != nil &&
		seed.Spec.Settings.Scheduling.Visible

	hasNoTaints := checkTolerations(seed, tolerations)

	result := isDeletionTimesampt && seed.Spec.Settings.Scheduling.Visible && isReady && hasNoTaints
	if !result {
		slog.Info("seed rejected",
			"name", seed.Name,
			"isDeletionTimestamp", isDeletionTimesampt,
			"isVisible", isVisible,
			"hasNoTaints", hasNoTaints,
			"isReady", isReady)
	}
	return result
}

func ToProviderRegions(seeds []gardener_types.Seed, tolerations config.TolerationsConfig) (out types.Providers) {
	defer LogWithDuration(time.Now(), "conversion complete")

	out = types.Providers{}
	for _, seed := range seeds {
		if seedCanBeUsed(&seed, tolerations) {
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
