package cli

import (
	"flag"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Gardener struct {
	KubeconfigPath   string
	Timeout          string
	SeedMapName      string
	SeedMapNamespace string
}

type Config struct {
	Gardener Gardener
	LogLevel string
}

func (c *Config) seedMapKey() client.ObjectKey {
	return client.ObjectKey{
		Namespace: c.Gardener.SeedMapNamespace,
		Name:      c.Gardener.SeedMapName,
	}
}

var ErrInvalidValue = fmt.Errorf("invalid value")

func validate[T any](value T, rulez []func(T) bool) error {
	for _, isValid := range rulez {
		if !isValid(value) {
			return fmt.Errorf("%w: %v", ErrInvalidValue, value)
		}
	}
	return nil
}

func isNotEmpty(s string) bool {
	return s != ""
}

func isValidDuration(s string) bool {
	_, err := time.ParseDuration(s)
	return err == nil
}

func isValidLogLevel(s string) bool {
	_, found := logLevelMapping[s]
	return found
}

func (c *Config) Validate() error {
	for _, item := range []struct {
		fieldValues []string
		validators  []func(string) bool
	}{
		{
			fieldValues: []string{
				c.Gardener.KubeconfigPath,
				c.Gardener.SeedMapName,
				c.Gardener.SeedMapNamespace,
			},
			validators: []func(string) bool{isNotEmpty},
		},
		{
			fieldValues: []string{
				c.Gardener.Timeout,
			},
			validators: []func(string) bool{isValidDuration},
		},
		{
			fieldValues: []string{
				c.LogLevel,
			},
			validators: []func(string) bool{isValidLogLevel},
		},
	} {
		for _, isValid := range item.validators {
			for _, value := range item.fieldValues {
				if !isValid(value) {
					return fmt.Errorf("%w: %v", ErrInvalidValue, value)
				}
			}
		}
	}

	return nil
}

const (
	FlagNameGardenerKubeconfigPath            = "gardener-kubeconfig-path"
	FlagNameGardenerSeedConfigMapName         = "gardener-seed-map-name"
	FlagNameGardenerSeedConfigMapNamespace    = "gardener-seed-map-namespace"
	FlagNameGardenerTimeout                   = "gardener-timeout"
	FlagDefaultGardenerKubeconfigPath         = "/gardener/kubeconfig"
	FlagDefaultGardenerSeedConfigMapName      = "gardener-seeds-cache"
	FlagDefaultGardenerSeedConfigMapNamespace = "kcp-system"
	FlagDefaultGardenerTimeout                = "10s"
	FlagNameLogLevel                          = "log-level"
	FlagDefaultLogLevel                       = "INFO"
)

func logLevelMappingKeys() []string {
	out := make([]string, 0, len(logLevelMapping))
	for key := range logLevelMapping {
		out = append(out, key)
	}
	return out
}

func NewConfigFromFlags() (Config, error) {
	out := Config{}

	flag.StringVar(&out.Gardener.KubeconfigPath, FlagNameGardenerKubeconfigPath, FlagDefaultGardenerKubeconfigPath, "A path to gardener kubeconfig file.")
	flag.StringVar(&out.Gardener.SeedMapName, FlagNameGardenerSeedConfigMapName, FlagDefaultGardenerSeedConfigMapName, "The name of the config-map that will store gardener seeds.")
	flag.StringVar(&out.Gardener.SeedMapNamespace, FlagNameGardenerSeedConfigMapNamespace, FlagDefaultGardenerSeedConfigMapNamespace, "The namespace of the config-map that will store gardener seeds.")
	flag.StringVar(&out.Gardener.Timeout, FlagNameGardenerTimeout, FlagDefaultGardenerTimeout, "Gardener client timeout duration.")
	flag.StringVar(&out.LogLevel, FlagNameLogLevel, FlagDefaultLogLevel, fmt.Sprintf("One of: %s", strings.Join(logLevelMappingKeys(), ",")))

	flag.Parse()

	if err := out.Validate(); err != nil {
		return Config{}, err
	}

	slog.Info("configuration parsed",
		FlagNameGardenerKubeconfigPath, out.Gardener.KubeconfigPath,
		FlagNameGardenerSeedConfigMapName, out.Gardener.SeedMapName,
		FlagNameGardenerSeedConfigMapNamespace, out.Gardener.SeedMapNamespace,
		FlagNameGardenerTimeout, out.Gardener.Timeout,
	)

	return out, nil
}
