package profile

import (
	"os"
	"time"

	"github.com/josh-tracey/scribe"
)

var profiling = os.Getenv("PROFILING")

func Duration(logger scribe.Logger, invocation time.Time, name string) {

	if profiling == "true" {
		elapsed := time.Since(invocation)
		logger.Info("%s: %s", name, elapsed)
	}
}
