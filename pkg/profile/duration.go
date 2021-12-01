package profile

import (
	"os"
	"time"

	"github.com/josh-tracey/eventual-agent/pkg/logging"
)

var profiling = os.Getenv("PROFILING")

func Duration(logger logging.Logger, invocation time.Time, name string) {

	if profiling == "true" {
		elapsed := time.Since(invocation)
		logger.Info("%s: %s", name, elapsed)
	}
}
