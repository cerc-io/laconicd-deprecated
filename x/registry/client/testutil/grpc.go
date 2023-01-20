package testutil

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cerc-io/laconicd/x/registry/client/cli"
	nstypes "github.com/cerc-io/laconicd/x/registry/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/rest"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

func (s *IntegrationTestSuite) TestGRPCQueryParams() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqURL := val.APIAddress + "/vulcanize/registry/v1beta1/params"

	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
	}{
		{
			"invalid url",
			reqURL + "/asdasd",
			true,
			"",
		},
		{
			"Success",
			reqURL,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, _ := rest.GetRequest(tc.url)
			require := s.Require()
			if tc.expectErr {
				require.Contains(string(resp), tc.errorMsg)
			} else {
				var response nstypes.QueryParamsResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				params := nstypes.DefaultParams()
				params.RecordRent = sdk.NewCoin(s.cfg.BondDenom, nstypes.DefaultRecordRent)
				params.RecordRentDuration = 10 * time.Second
				params.AuthorityGracePeriod = 10 * time.Second
				sr.Equal(response.GetParams().String(), params.String())
			}
		})
	}
}

//nolint:all
func (s *IntegrationTestSuite) TestGRPCQueryWhoIs() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqUrl := val.APIAddress + "/vulcanize/registry/v1beta1/whois/%s"
	authorityName := "QueryWhoIS"
	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
		preRun    func(authorityName string)
	}{
		{
			"invalid url",
			reqUrl + "/asdasd",
			true,
			"",
			func(authorityName string) {
			},
		},
		{
			"Success",
			reqUrl,
			false,
			"",
			func(authorityName string) {
				clientCtx := val.ClientCtx
				cmd := cli.GetCmdReserveName()
				args := []string{
					authorityName,
					fmt.Sprintf("--owner=%s", accountAddress),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expectErr {
				tc.preRun(authorityName)
				tc.url = fmt.Sprintf(tc.url, authorityName)
			}
			resp, _ := rest.GetRequest(tc.url)
			require := s.Require()
			if tc.expectErr {
				require.Contains(string(resp), tc.errorMsg)
			} else {
				var response nstypes.QueryWhoisResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				sr.Equal(nstypes.AuthorityActive, response.GetNameAuthority().Status)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCQueryLookup() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqURL := val.APIAddress + "/vulcanize/registry/v1beta1/lookup?crn=%s"
	authorityName := "QueryLookUp"

	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
		preRun    func(authorityName string)
	}{
		{
			"invalid url",
			reqURL + "/asdasd",
			true,
			"",
			func(authorityName string) {
			},
		},
		{
			"Success",
			reqURL,
			false,
			"",
			func(authorityName string) {
				// create name record
				createNameRecord(authorityName, s)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expectErr {
				tc.preRun(authorityName)
				tc.url = fmt.Sprintf(reqURL, fmt.Sprintf("crn://%s/", authorityName))
			}
			resp, _ := rest.GetRequest(tc.url)
			if tc.expectErr {
				sr.Contains(string(resp), tc.errorMsg)
			} else {
				var response nstypes.QueryLookupCrnResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				sr.NotZero(len(response.Name.Latest.Id))
			}
		})
	}
}

//nolint:all
func (s *IntegrationTestSuite) TestGRPCQueryRecordExpiryQueue() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqUrl := val.APIAddress + "/vulcanize/registry/v1beta1/record-expiry"

	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
		preRun    func(bondId string)
	}{
		{
			"invalid url",
			reqUrl + "/asdasd",
			true,
			"",
			func(bondId string) {
			},
		},
		{
			"Success",
			reqUrl,
			false,
			"",
			func(bondId string) {
				dir, err := os.Getwd()
				sr.NoError(err)
				payloadPath := dir + "/../../helpers/examples/service_provider_example.yml"
				args := []string{
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}
				args = append([]string{payloadPath, bondId}, args...)
				clientCtx := val.ClientCtx
				cmd := cli.GetCmdSetRecord()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expectErr {
				tc.preRun(s.bondID)
			}
			// wait 12 seconds for records expires
			time.Sleep(time.Second * 12)
			resp, _ := rest.GetRequest(tc.url)
			require := s.Require()
			if tc.expectErr {
				require.Contains(string(resp), tc.errorMsg)
			} else {
				var response nstypes.QueryGetRecordExpiryQueueResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				sr.NotZero(len(response.GetRecords()))
			}
		})
	}
}

//nolint:all
func (s *IntegrationTestSuite) TestGRPCQueryAuthorityExpiryQueue() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqUrl := val.APIAddress + "/vulcanize/registry/v1beta1/authority-expiry"

	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
		preRun    func(authorityName string)
	}{
		{
			"invalid url",
			reqUrl + "/asdasd",
			true,
			"",
			func(authorityName string) {
			},
		},
		{
			"Success",
			reqUrl,
			false,
			"",
			func(authorityName string) {
				// reserving the name
				clientCtx := val.ClientCtx
				cmd := cli.GetCmdReserveName()
				args := []string{
					authorityName,
					fmt.Sprintf("--owner=%s", accountAddress),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expectErr {
				tc.preRun("QueryAuthorityExpiryQueue")
			}
			// wait 12 seconds to name authorites expires
			time.Sleep(time.Second * 12)

			resp, _ := rest.GetRequest(tc.url)
			require := s.Require()
			if tc.expectErr {
				require.Contains(string(resp), tc.errorMsg)
			} else {
				var response nstypes.QueryGetAuthorityExpiryQueueResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				// removed from expiry queue as no bond set
				sr.Zero(len(response.GetAuthorities()))
			}
		})
	}
}

//nolint:all
func (s *IntegrationTestSuite) TestGRPCQueryListRecords() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqUrl := val.APIAddress + "/vulcanize/registry/v1beta1/records"

	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
		preRun    func(bondId string)
	}{
		{
			"invalid url",
			reqUrl + "/asdasd",
			true,
			"",
			func(bondId string) {
			},
		},
		{
			"Success",
			reqUrl,
			false,
			"",
			func(bondId string) {
				dir, err := os.Getwd()
				sr.NoError(err)
				payloadPath := dir + "/../../helpers/examples/service_provider_example.yml"
				args := []string{
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}
				args = append([]string{payloadPath, bondId}, args...)
				clientCtx := val.ClientCtx
				cmd := cli.GetCmdSetRecord()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expectErr {
				tc.preRun(s.bondID)
			}
			resp, _ := rest.GetRequest(tc.url)
			require := s.Require()
			if tc.expectErr {
				require.Contains(string(resp), tc.errorMsg)
			} else {
				var response nstypes.QueryListRecordsResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				sr.NotZero(len(response.GetRecords()))
				sr.Equal(s.bondID, response.GetRecords()[0].GetBondId())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCQueryGetRecordByID() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqURL := val.APIAddress + "/vulcanize/registry/v1beta1/records/%s"

	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
		preRun    func(bondId string) string
	}{
		{
			"invalid url",
			reqURL + "/asdasd",
			true,
			"",
			func(bondId string) string {
				return ""
			},
		},
		{
			"Success",
			reqURL,
			false,
			"",
			func(bondId string) string {
				// creating the record
				createRecord(bondId, s)

				// list the records
				clientCtx := val.ClientCtx
				cmd := cli.GetCmdList()
				args := []string{
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				}
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var records []nstypes.RecordType
				err = json.Unmarshal(out.Bytes(), &records)
				sr.NoError(err)
				return records[0].ID
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var recordID string
			if !tc.expectErr {
				recordID = tc.preRun(s.bondID)
				tc.url = fmt.Sprintf(reqURL, recordID)
			}
			resp, _ := rest.GetRequest(tc.url)
			require := s.Require()
			if tc.expectErr {
				require.Contains(string(resp), tc.errorMsg)
			} else {
				var response nstypes.QueryRecordByIDResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				record := response.GetRecord()
				sr.NotZero(len(record.GetId()))
				sr.Equal(record.GetId(), recordID)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCQueryGetRecordByBondID() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqURL := val.APIAddress + "/vulcanize/registry/v1beta1/records-by-bond-id/%s"

	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
		preRun    func(bondId string)
	}{
		{
			"invalid url",
			reqURL + "/asdasd",
			true,
			"",
			func(bondId string) {
			},
		},
		{
			"Success",
			reqURL,
			false,
			"",
			func(bondId string) {
				// creating the record
				createRecord(bondId, s)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expectErr {
				tc.preRun(s.bondID)
				tc.url = fmt.Sprintf(reqURL, s.bondID)
			}
			resp, _ := rest.GetRequest(tc.url)
			require := s.Require()
			if tc.expectErr {
				require.Contains(string(resp), tc.errorMsg)
			} else {
				var response nstypes.QueryRecordByBondIDResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				records := response.GetRecords()
				sr.NotZero(len(records))
				sr.Equal(records[0].GetBondId(), s.bondID)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCQueryGetRegistryModuleBalance() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqURL := val.APIAddress + "/vulcanize/registry/v1beta1/balance"

	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
		preRun    func(bondId string)
	}{
		{
			"invalid url",
			reqURL + "/asdasd",
			true,
			"",
			func(bondId string) {
			},
		},
		{
			"Success",
			reqURL,
			false,
			"",
			func(bondId string) {
				// creating the record
				createRecord(bondId, s)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expectErr {
				tc.preRun(s.bondID)
			}
			resp, _ := rest.GetRequest(tc.url)
			require := s.Require()
			if tc.expectErr {
				require.Contains(string(resp), tc.errorMsg)
			} else {
				var response nstypes.GetRegistryModuleBalanceResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				sr.NotZero(len(response.GetBalances()))
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGRPCQueryNamesList() {
	val := s.network.Validators[0]
	sr := s.Require()
	reqURL := val.APIAddress + "/vulcanize/registry/v1beta1/names"

	testCases := []struct {
		name      string
		url       string
		expectErr bool
		errorMsg  string
		preRun    func(authorityName string)
	}{
		{
			"invalid url",
			reqURL + "/asdasd",
			true,
			"",
			func(authorityName string) {
			},
		},
		{
			"Success",
			reqURL,
			false,
			"",
			func(authorityName string) {
				// create name record
				createNameRecord(authorityName, s)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expectErr {
				tc.preRun("ListNameRecords")
			}
			resp, _ := rest.GetRequest(tc.url)
			require := s.Require()
			if tc.expectErr {
				require.Contains(string(resp), tc.errorMsg)
			} else {
				var response nstypes.QueryListNameRecordsResponse
				err := val.ClientCtx.Codec.UnmarshalJSON(resp, &response)
				sr.NoError(err)
				sr.NotZero(len(response.GetNames()))
			}
		})
	}
}

func createRecord(bondID string, s *IntegrationTestSuite) {
	val := s.network.Validators[0]
	sr := s.Require()

	dir, err := os.Getwd()
	sr.NoError(err)
	payloadPath := dir + "/../../helpers/examples/service_provider_example.yml"
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
	}
	args = append([]string{payloadPath, bondID}, args...)
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdSetRecord()

	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	sr.NoError(err)
	var d sdk.TxResponse
	err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
	sr.NoError(err)
	sr.Zero(d.Code)
}
