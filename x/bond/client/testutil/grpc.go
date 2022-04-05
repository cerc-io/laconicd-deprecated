package testutil

import (
	"fmt"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/types/rest"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tharsis/ethermint/x/bond/client/cli"
	bondtypes "github.com/tharsis/ethermint/x/bond/types"
)

func (s *IntegrationTestSuite) TestGRPCGetBonds() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqUrl := fmt.Sprintf("%s/vulcanize/bond/v1beta1/bonds", val.APIAddress)

	testCases := []struct {
		name     string
		url      string
		expErr   bool
		errorMsg string
		preRun   func()
	}{
		{
			"invalid request with headers",
			reqUrl + "asdasdas",
			true,
			"",
			func() {

			},
		},
		{
			"valid request",
			reqUrl,
			false,
			"",
			func() {
				s.createBond()
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, _ := rest.GetRequest(tc.url)
			if tc.expErr {
				sr.Contains(string(resp), tc.errorMsg)
			} else {
				var response bondtypes.QueryGetBondsResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				sr.NotZero(len(response.GetBonds()))
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCGetParams() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqUrl := fmt.Sprintf("%s/vulcanize/bond/v1beta1/params", val.APIAddress)

	resp, err := rest.GetRequest(reqUrl)
	s.Require().NoError(err)

	var params bondtypes.QueryParamsResponse
	err = val.ClientCtx.Codec.UnmarshalJSON(resp, &params)

	sr.NoError(err)
	sr.Equal(params.GetParams().MaxBondAmount, bondtypes.DefaultParams().MaxBondAmount)
}

func (s *IntegrationTestSuite) TestGRPCGetBondsByOwner() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqUrl := val.APIAddress + "/vulcanize/bond/v1beta1/by-owner/%s"

	testCases := []struct {
		name   string
		url    string
		expErr bool
		preRun func()
	}{
		{
			"empty list",
			fmt.Sprintf(reqUrl, "asdasd"),
			true,
			func() {

			},
		},
		{
			"valid request",
			fmt.Sprintf(reqUrl, s.accountAddress),
			false,
			func() {
				s.createBond()
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expErr {
				tc.preRun()
			}

			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			var bonds bondtypes.QueryGetBondsByOwnerResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bonds)
			sr.NoError(err)
			if tc.expErr {
				sr.Empty(bonds.GetBonds())
			} else {
				bondsList := bonds.GetBonds()
				sr.NotZero(len(bondsList))
				sr.Equal(s.accountAddress, bondsList[0].GetOwner())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCGetBondById() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqUrl := val.APIAddress + "/vulcanize/bond/v1beta1/bonds/%s"

	testCases := []struct {
		name   string
		url    string
		expErr bool
		preRun func() string
	}{
		{
			"invalid request",
			fmt.Sprintf(reqUrl, "asdadad"),
			true,
			func() string {
				return ""
			},
		},
		{
			"valid request",
			reqUrl,
			false,
			func() string {
				// creating the bond
				s.createBond()

				// getting the bonds list and returning the bond-id
				clientCtx := val.ClientCtx
				cmd := cli.GetQueryBondLists()
				args := []string{
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				}
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var queryResponse bondtypes.QueryGetBondsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				sr.NoError(err)

				// extract bond id from bonds list
				bond := queryResponse.GetBonds()[0]
				return bond.GetId()
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var bondId string
			if !tc.expErr {
				bondId = tc.preRun()
				tc.url = fmt.Sprintf(reqUrl, bondId)
			}

			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			var bonds bondtypes.QueryGetBondByIdResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bonds)

			if tc.expErr {
				sr.Empty(bonds.GetBond().GetId())
			} else {
				sr.NoError(err)
				sr.NotZero(bonds.GetBond().GetId())
				sr.Equal(bonds.GetBond().GetId(), bondId)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCGetBondModuleBalance() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqUrl := fmt.Sprintf("%s/vulcanize/bond/v1beta1/balance", val.APIAddress)

	// creating the bond
	s.createBond()

	s.Run("valid request", func() {
		resp, err := rest.GetRequest(reqUrl)
		sr.NoError(err)

		var response bondtypes.QueryGetBondModuleBalanceResponse
		err = val.ClientCtx.Codec.UnmarshalJSON(resp, &response)

		sr.NoError(err)
		sr.False(response.GetBalance().IsZero())
	})
}
