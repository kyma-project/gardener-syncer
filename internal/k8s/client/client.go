package client

import (
	"log/slog"

	"github.com/kyma-project/infrastructure-manager/pkg/gardener"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Options struct {
	KubeconfigPath        string
	AdditionalAddToSchema []func(*runtime.Scheme) error
}

func New(opt Options) (k8sClient client.Client, err error) {

	scheme := runtime.NewScheme()
	for _, register := range opt.AdditionalAddToSchema {
		if err := register(scheme); err != nil {
			return nil, err
		}
	}

	slog.Info("schema registered")

	getRestConfig := config.GetConfig
	if opt.KubeconfigPath != "" {
		getRestConfig = func() (*rest.Config, error) {
			return gardener.NewRestConfigFromFile(opt.KubeconfigPath)
		}
	}

	restConfig, err := getRestConfig()
	if err != nil {
		return nil, err
	}

	gardenerClient, err := client.New(restConfig, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}

	return gardenerClient, nil
}
