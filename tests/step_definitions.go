package tests

import (
	"context"
	"log"

	"github.com/cucumber/godog"
)

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {})
}

func InitializeScenarios(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		log.Printf("Initializing API for scenario: %s\n", sc.Name)

		return ctx, nil
	})

	ctx.Given(`^the API is running`, givenAPIRunning)

	ctx.When(`^I receive "([^"]*)" request to "([^"]*)" with "(.*)" body`, whenReceiveRequest)

	ctx.Then(`the response http status code should be "([^"]*)"`, thenResponseHTTPStatusCodeShouldBe)

	ctx.Then(`^the response body should be "(.*)"`, thenResponseBodyShouldBe)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, _ error) (context.Context, error) {
		log.Printf("Stopping API for scenario: %s\n", sc.Name)

		return ctx, tearDownScenario(ctx)
	})
}
