package seeker

import (
	"context"
	"github.com/kyma-project/infrastructure-manager/pkg/config"
	"time"

	gardener_types "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/kyma-project/gardener-syncer/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type List func(context.Context, client.ObjectList, ...client.ListOption) error

type FetchSeeds func() (types.Providers, error)

type FetchSeedsOpts struct {
	Timeout time.Duration
	List
}

func BuildFetchSeedFn(opts FetchSeedsOpts, tolerations config.TolerationsConfig) FetchSeeds {
	return func() (types.Providers, error) {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer LogWithDuration(time.Now(), "fetching gardener data complete")
		defer cancel()

		seeds, err := listSeeds(ctx, opts.List)
		if err != nil {
			return nil, err
		}

		return ToProviderRegions(seeds.Items, tolerations), nil
	}
}

func listSeeds(ctx context.Context, list List) (seeds gardener_types.SeedList, err error) {
	defer func() {
		LogWithDuration(time.Now(), "gardener-seed list complete", "count", len(seeds.Items))
	}()

	if err = list(ctx, &seeds); err != nil {
		return gardener_types.SeedList{}, err
	}

	return seeds, nil
}
