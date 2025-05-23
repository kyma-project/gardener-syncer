package seeker

import (
	"context"
	"time"

	log "log/slog"

	"github.com/kyma-project/gardener-syncer/pkg/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var FieldManagerName = "gardener-syncer"

type Patch func(context.Context, client.Object, client.Patch, ...client.PatchOption) error

type Get func(context.Context, client.ObjectKey, client.Object, ...client.GetOption) error

type Store func(types.Providers) error

type StoreOpts struct {
	Timeout time.Duration
	Key     client.ObjectKey
	Patch
	Get
	Convert[types.Providers, map[string]string]
}

func logWithDuration(startTime time.Time) {
	duration := time.Now().Sub(startTime)
	log.With("duration", duration).Info("done")
}

func BuildStoreFn(opts StoreOpts) Store {
	return func(data types.Providers) (err error) {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()

		var cm corev1.ConfigMap
		fetch := func() error {
			log.With("config-map", opts.Key).Info("fetching")
			defer logWithDuration(time.Now())
			return opts.Get(ctx, opts.Key, &cm)
		}

		if err = fetch(); err != nil && !errors.IsNotFound(err) {
			return err
		}

		cm.Name = opts.Key.Name
		cm.Namespace = opts.Key.Namespace
		cm.Data, err = opts.Convert(data)
		cm.TypeMeta.Kind = "ConfigMap"
		cm.TypeMeta.APIVersion = "v1"
		cm.ManagedFields = nil

		if err != nil {
			return err
		}

		return opts.Patch(ctx, &cm, client.Apply, &client.PatchOptions{
			FieldManager: FieldManagerName,
		})
	}
}
