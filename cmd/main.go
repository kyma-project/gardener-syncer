package main

import (
	log "log/slog"
	"net/http"
	"net/url"
	"os"

	cli "github.com/kyma-project/gardener-syncer/internal"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Error("application failed", "error", err.Error())
		haltIstioSidecar() // os.Exit is not called in defer, so we need to call it here as well
		os.Exit(1)
	}
}
