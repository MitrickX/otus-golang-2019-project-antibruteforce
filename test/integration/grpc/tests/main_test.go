package tests

import (
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
)

var runnerOptions = godog.Options{
	Format:      "pretty", // progress, pretty
	Paths:       []string{"../features/"},
	Randomize:   0,
	Concurrency: 0,
}

// Test entry point
func TestMain(m *testing.M) {
	config := GetConfig()

	if len(config.RunnerPaths) > 0 {
		runnerOptions.Paths = config.RunnerPaths
	}

	status := godog.RunWithOptions("integration", func(s *godog.Suite) {
		t := newFeatureTest()
		FeatureContext(s, t)
	}, runnerOptions)

	os.Exit(status)
}
