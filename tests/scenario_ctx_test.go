package tests

import (
	"github.com/jrogala/spotify-cli/client"
)

// scenarioCtx holds per-scenario state.
type scenarioCtx struct {
	mock   *MockServer
	client *client.Client
	err    error
	result any
}

func newScenarioCtx(mock *MockServer) *scenarioCtx {
	return &scenarioCtx{
		mock:   mock,
		client: client.NewWithBaseURL("test-token", mock.BaseURL()),
	}
}
