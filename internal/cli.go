package cli

import (
	"fmt"
	"log/slog"
	log "log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/kyma-project/gardener-syncer/internal/k8s/client"
	seeker "github.com/kyma-project/gardener-syncer/pkg"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	defaultKcpClientTimeout = time.Second * 10
	logLevelMapping         = map[string]log.Level{
		"INFO":  log.LevelInfo,
		"DEBUG": log.LevelDebug,
	}
)

func Run() error {
	defer seeker.LogWithDuration(time.Now(), "application finished")
	defer haltIstioSidecar()

	cfg, err := NewConfigFromFlags()
	if err != nil {
		return err
	}

	logLevel := mustParseLogLevel(cfg.LogLevel)
	slog.SetLogLoggerLevel(logLevel)

	kcpClient, err := client.New(client.Options{
		AdditionalAddToSchema: []func(*runtime.Scheme) error{
			corev1.AddToScheme,
		},
	}, "kcp")

	if err != nil {
		return err
	}

	store := seeker.BuildStoreFn(seeker.StoreOpts{
		Key:     cfg.seedMapKey(),
		Patch:   kcpClient.Patch,
		Get:     kcpClient.Get,
		Convert: seeker.ToConfigMap,
		Timeout: defaultKcpClientTimeout,
	})

	gardenerClient, err := client.New(
		client.Options{
			KubeconfigPath: cfg.Gardener.KubeconfigPath,
			AdditionalAddToSchema: []func(*runtime.Scheme) error{
				v1beta1.AddToScheme,
			},
		},
		"gardener",
	)

	if err != nil {
		return err
	}

	gardenerTimeout := mustParseDuration(cfg.Gardener.Timeout)
	fetch := seeker.BuildFetchSeedFn(seeker.FetchSeedsOpts{
		List:    gardenerClient.List,
		Timeout: gardenerTimeout,
	})

	sync := seeker.BuildSyncFn(store, fetch)
	return sync()
}

func mustParseDuration(s string) time.Duration {
	out, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("invalid duration value: %s", s))
	}
	return out
}

func mustParseLogLevel(s string) log.Level {
	level, found := logLevelMapping[s]
	if !found {
		panic(fmt.Sprintf("invalid log level: %s", s))
	}
	return level
}

func haltIstioSidecar() {
	log.Info("# HALT ISTIO SIDECAR #")
	resp, err := http.PostForm("http://127.0.0.1:15020/quitquitquit", url.Values{})

	if err != nil {
		log.Error("unable to send post request to quit Istio sidecar", "error", err)
		return
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Info("Stopping istio sidecar, ", "response status", resp.StatusCode)
		return
	}
}
