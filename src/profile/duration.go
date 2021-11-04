package profile

import (
	"os"
	"time"

	"adriftdev.com/eventual-agent/src/logging"
)

var profiling = os.Getenv("PROFILING")

func Duration(logger logging.Logger, invocation time.Time, name string) {

	if profiling == "true" {
		elapsed := time.Since(invocation)
		logger.Info("%s: %s", name, elapsed)
	}
}
