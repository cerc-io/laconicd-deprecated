package testutil

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tharsis/ethermint/testutil/network"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = 2
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
