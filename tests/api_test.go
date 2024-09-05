package tests

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestAPIFeature(t *testing.T) {
	suite := godog.TestSuite{
		Name:                 "API User",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenarios,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"./api.feature"},
			TestingT: t,
			Strict:   true,
		},
	}

	if suite.Run() != 0 {
		t.Fail()
	}
}
