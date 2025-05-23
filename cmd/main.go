package main

import (
	log "log/slog"
	"os"

	cli "github.com/kyma-project/gardener-syncer/internal"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Error("application failed", "error", err.Error())
		os.Exit(1)
	}
}
