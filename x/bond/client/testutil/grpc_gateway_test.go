package testutil

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/testutil"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/rest"
	bondtypes "github.com/tharsis/ethermint/x/bond/types"
)

func (s *IntegrationTestSuite) TestGetBondsGRPC() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name      string
		url       string
		headers   map[string]string
		noOfBonds int
		expErr    bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/bond/v1beta1/bonds", val.APIAddress),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/bond/v1beta1/bonds", val.APIAddress),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var bonds bondtypes.QueryGetBondsResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bonds)

			if tc.expErr {
				sr.Empty(bonds.GetBonds())
			} else {
				sr.NoError(err)
				sr.Len(bonds.GetBonds(), tc.noOfBonds)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetParamsGRPC() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name      string
		url       string
		headers   map[string]string
		noOfBonds int
		expErr    bool
	}{
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/bond/v1beta1/params", val.APIAddress),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			var params bondtypes.QueryParamsResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &params)

			if tc.expErr {
				sr.Empty(params.GetParams())
			} else {
				sr.NoError(err)
				sr.Equal(params.GetParams().MaxBondAmount, bondtypes.DefaultParams().MaxBondAmount)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetBondsByOwnerGRPC() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name      string
		url       string
		headers   map[string]string
		noOfBonds int
		expErr    bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/bond/v1beta1/by-owner", val.APIAddress),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/bond/v1beta1/by-owner/%s,", val.APIAddress, val.Address.String()),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var bonds bondtypes.QueryGetBondsByOwnerResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bonds)

			if tc.expErr {
				sr.Empty(bonds.GetBonds())
			} else {
				sr.NoError(err)
				sr.Len(bonds.GetBonds(), tc.noOfBonds)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetBondByIdGRPC() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name      string
		url       string
		headers   map[string]string
		noOfBonds int
		expErr    bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/bond/v1beta1/bonds/%s", val.APIAddress, "asdadad"),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/bond/v1beta1/bonds/%s,", val.APIAddress, "asdadad"),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var bonds bondtypes.QueryGetBondByIdResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bonds)

			if tc.expErr {
				sr.Empty(bonds.GetBond().GetId())
			} else {
				sr.NoError(err)
				sr.Len(bonds.GetBond().GetId(), tc.noOfBonds)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetBondModuleBalanceGRPC() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCases := []struct {
		name      string
		url       string
		headers   map[string]string
		noOfBonds int
		expErr    bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/bond/v1beta1/balance", val.APIAddress),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/bond/v1beta1/balance", val.APIAddress),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var bonds bondtypes.QueryGetBondModuleBalanceResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bonds)

			if tc.expErr {
				sr.Empty(bonds.GetBalance())
			} else {
				sr.NoError(err)
				sr.True(bonds.GetBalance().IsZero())
			}
		})
	}
}
