package main

import (
	log "log/slog"
	"net/http"
	"net/url"
	"os"

	cli "github.com/kyma-project/gardener-syncer/internal"
)

func main() {
	defer haltIstioSidecar()

	if err := cli.Run(); err != nil {
		log.Error("application failed", "error", err.Error())
		os.Exit(1)
	}
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
