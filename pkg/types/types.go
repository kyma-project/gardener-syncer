package types

import (
	"log/slog"
	"slices"
)

type ProviderInfo struct {
	SeedRegions []string `json:"seedRegions"`
}

type Providers map[string]ProviderInfo

func (s *Providers) Add(provider, regionName string) {
	providerInfo := (*s)[provider]
	if slices.Contains(providerInfo.SeedRegions, regionName) {
		slog.Debug("region name already collected", "provider", provider, "regionName", regionName)
		return
	}

	providerInfo.SeedRegions = append(providerInfo.SeedRegions, regionName)
	(*s)[provider] = providerInfo
}
