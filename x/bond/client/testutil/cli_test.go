package testutil

import (
	"github.com/stretchr/testify/suite"
	"github.com/tharsis/ethermint/testutil/network"
	"testing"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = 1
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
