package seeker

import (
	"time"
)

type Sync func() error

func BuildSyncFn(store Store, fetch FetchSeeds) Sync {
	return func() (err error) {
		defer LogWithDuration(time.Now(), "synchronisation complete")

		providerRegions, err := fetch()
		if err != nil {
			return err
		}

		if err := store(providerRegions); err != nil {
			return err
		}

		return nil
	}
}
