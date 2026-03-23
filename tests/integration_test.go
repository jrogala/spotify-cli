package tests

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	mock := NewMockServer()
	defer mock.Close()

	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			ctx := newScenarioCtx(mock)
			registerSteps(sc, ctx)

			sc.Before(func(goCtx context.Context, sc2 *godog.Scenario) (context.Context, error) {
				mock.Reset()
				return goCtx, nil
			})
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		os.Exit(1)
	}
}
