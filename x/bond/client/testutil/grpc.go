package testutil

import (
	"fmt"

	"github.com/cerc-io/laconicd/x/bond/client/cli"
	bondtypes "github.com/cerc-io/laconicd/x/bond/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/rest"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

func (s *IntegrationTestSuite) TestGRPCGetBonds() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqURL := fmt.Sprintf("%s/vulcanize/bond/v1beta1/bonds", val.APIAddress)

	testCases := []struct {
		name     string
		url      string
		expErr   bool
		errorMsg string
		preRun   func()
	}{
		{
			"invalid request with headers",
			reqURL + "asdasdas",
			true,
			"",
			func() {
			},
		},
		{
			"valid request",
			reqURL,
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
	reqURL := fmt.Sprintf("%s/vulcanize/bond/v1beta1/params", val.APIAddress)

	resp, err := rest.GetRequest(reqURL)
	s.Require().NoError(err)

	var params bondtypes.QueryParamsResponse
	err = val.ClientCtx.Codec.UnmarshalJSON(resp, &params)

	sr.NoError(err)
	sr.Equal(params.GetParams().MaxBondAmount, bondtypes.DefaultParams().MaxBondAmount)
}

func (s *IntegrationTestSuite) TestGRPCGetBondsByOwner() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqURL := val.APIAddress + "/vulcanize/bond/v1beta1/by-owner/%s"

	testCases := []struct {
		name   string
		url    string
		expErr bool
		preRun func()
	}{
		{
			"empty list",
			fmt.Sprintf(reqURL, "asdasd"),
			true,
			func() {
			},
		},
		{
			"valid request",
			fmt.Sprintf(reqURL, s.accountAddress),
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

func (s *IntegrationTestSuite) TestGRPCGetBondByID() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqURL := val.APIAddress + "/vulcanize/bond/v1beta1/bonds/%s"

	testCases := []struct {
		name   string
		url    string
		expErr bool
		preRun func() string
	}{
		{
			"invalid request",
			fmt.Sprintf(reqURL, "asdadad"),
			true,
			func() string {
				return ""
			},
		},
		{
			"valid request",
			reqURL,
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
				return bond.GetID()
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var bondID string
			if !tc.expErr {
				bondID = tc.preRun()
				tc.url = fmt.Sprintf(reqURL, bondID)
			}

			resp, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)

			var bonds bondtypes.QueryGetBondByIDResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bonds)

			if tc.expErr {
				sr.Empty(bonds.GetBond().GetID())
			} else {
				sr.NoError(err)
				sr.NotZero(bonds.GetBond().GetID())
				sr.Equal(bonds.GetBond().GetID(), bondID)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCGetBondModuleBalance() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqURL := fmt.Sprintf("%s/vulcanize/bond/v1beta1/balance", val.APIAddress)

	// creating the bond
	s.createBond()

	s.Run("valid request", func() {
		resp, err := rest.GetRequest(reqURL)
		sr.NoError(err)

		var response bondtypes.QueryGetBondModuleBalanceResponse
		err = val.ClientCtx.Codec.UnmarshalJSON(resp, &response)

		sr.NoError(err)
		sr.False(response.GetBalance().IsZero())
	})
}
