package keeper_test

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cerc-io/laconicd/x/registry/client/cli"
	"github.com/cerc-io/laconicd/x/registry/helpers"
	"github.com/cerc-io/laconicd/x/registry/keeper"
	registrytypes "github.com/cerc-io/laconicd/x/registry/types"
)

func (suite *KeeperTestSuite) TestGrpcQueryParams() {
	grpcClient := suite.queryClient

	testCases := []struct {
		msg string
		req *registrytypes.QueryParamsRequest
	}{
		{
			"Get Params",
			&registrytypes.QueryParamsRequest{},
		},
	}
	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			resp, _ := grpcClient.Params(context.Background(), test.req)
			defaultParams := registrytypes.DefaultParams()
			suite.Require().Equal(defaultParams.String(), resp.GetParams().String())
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcGetRecordLists() {
	grpcClient, ctx := suite.queryClient, suite.ctx
	sr := suite.Require()
	var recordId string
	examples := []string{
		"/../helpers/examples/binary_example.yml",
		"/../helpers/examples/chain_registration_example.yml",
		"/../helpers/examples/docker_image_example.yml",
		"/../helpers/examples/git_repo_example.yml",
		"/../helpers/examples/js_package_example.yml",
		"/../helpers/examples/responder_contract_example.yml",
		"/../helpers/examples/service_provider_example.yml",
		"/../helpers/examples/watcher_registration_example.yml",
		"/../helpers/examples/website_registration_example.yml",
	}
	testCases := []struct {
		msg           string
		req           *registrytypes.QueryListRecordsRequest
		createRecords bool
		expErr        bool
		noOfRecords   int
	}{
		{
			"Empty Records",
			&registrytypes.QueryListRecordsRequest{},
			false,
			false,
			0,
		},
		{
			"List Records",
			&registrytypes.QueryListRecordsRequest{},
			true,
			false,
			len(examples),
		},
		{
			"Filter with type",
			&registrytypes.QueryListRecordsRequest{
				Attributes: []*registrytypes.QueryListRecordsRequest_KeyValueInput{
					{
						Key: "type---",
						Value: &registrytypes.QueryListRecordsRequest_ValueInput{
							Type:    "string",
							String_: "WebsiteRegistrationRecord",
						},
					},
				},
				All: true,
			},
			true,
			false,
			1,
		},
		{
			"Filter with attributes ServiceProviderRecord",
			&registrytypes.QueryListRecordsRequest{
				Attributes: []*registrytypes.QueryListRecordsRequest_KeyValueInput{
					{
						Key: "x500---state_name---",
						Value: &registrytypes.QueryListRecordsRequest_ValueInput{
							Type:    "string",
							String_: "california",
						},
					},
				},
				All: true,
			},
			true,
			false,
			1,
		},
		{
			"Filter with attributes WatcherRegistrationRecord",
			&registrytypes.QueryListRecordsRequest{
				Attributes: []*registrytypes.QueryListRecordsRequest_KeyValueInput{
					{
						Key: "metadata---chain_reference---/---",
						Value: &registrytypes.QueryListRecordsRequest_ValueInput{
							Type:    "string",
							String_: "QmP8jTG1m9GSDJLCbeWhVSVgEzCPPwXRdCRuJtQ5Tz9Kc9",
						},
					},
				},
				All: true,
			},
			true,
			false,
			1,
		},
		{
			"Filter with attributes version",
			&registrytypes.QueryListRecordsRequest{
				Attributes: []*registrytypes.QueryListRecordsRequest_KeyValueInput{
					{
						Key: "version---",
						Value: &registrytypes.QueryListRecordsRequest_ValueInput{
							Type:    "string",
							String_: "1.0.0",
						},
					},
				},
				All: true,
			},
			true,
			false,
			7,
		},
	}
	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			if test.createRecords {
				for _, example := range examples {
					dir, err := os.Getwd()
					sr.NoError(err)
					payloadType, err := cli.GetPayloadFromFile(fmt.Sprint(dir, example))
					sr.NoError(err)
					payload, err := payloadType.ToPayload()
					sr.NoError(err)
					record, err := suite.app.RegistryKeeper.ProcessSetRecord(ctx, registrytypes.MsgSetRecord{
						BondId:  suite.bond.GetId(),
						Signer:  suite.accounts[0].String(),
						Payload: payload,
					})
					sr.NoError(err)
					sr.NotNil(record.ID)
				}
			}
			resp, err := grpcClient.ListRecords(context.Background(), test.req)
			if test.expErr {
				suite.Error(err)
			} else {
				sr.NoError(err)
				sr.Equal(test.noOfRecords, len(resp.GetRecords()))
				if test.createRecords {
					recordId = resp.GetRecords()[0].GetId()
					sr.NotZero(resp.GetRecords())
					sr.Equal(resp.GetRecords()[0].GetBondId(), suite.bond.GetId())

					for _, record := range resp.GetRecords() {
						bz, err := registrytypes.GetJSONBytesFromAny(*record.Attributes)
						sr.NoError(err)
						recAttr := helpers.UnMarshalMapFromJSONBytes(bz)
						for _, attr := range test.req.GetAttributes() {
							sr.Equal(keeper.GetAttributeValue(attr.Value), getAttributesFromRecord(recAttr, attr.Key))
						}
					}
				}
			}
		})
	}

	// Get the records by record id
	testCases1 := []struct {
		msg          string
		req          *registrytypes.QueryRecordByIDRequest
		createRecord bool
		expErr       bool
		noOfRecords  int
	}{
		{
			"Invalid Request without record id",
			&registrytypes.QueryRecordByIDRequest{},
			false,
			true,
			0,
		},
		{
			"With Record ID",
			&registrytypes.QueryRecordByIDRequest{
				Id: recordId,
			},
			true,
			false,
			1,
		},
	}
	for _, test := range testCases1 {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			resp, err := grpcClient.GetRecord(context.Background(), test.req)
			if test.expErr {
				suite.Error(err)
			} else {
				sr.NoError(err)
				sr.NotNil(resp.GetRecord())
				if test.createRecord {
					sr.Equal(resp.GetRecord().BondId, suite.bond.GetId())
					sr.Equal(resp.GetRecord().Id, recordId)
				}
			}
		})
	}

	// Get the records by record id
	testCasesByBondID := []struct {
		msg          string
		req          *registrytypes.QueryRecordByBondIDRequest
		createRecord bool
		expErr       bool
		noOfRecords  int
	}{
		{
			"Invalid Request without bond id",
			&registrytypes.QueryRecordByBondIDRequest{},
			false,
			true,
			0,
		},
		{
			"With Bond ID",
			&registrytypes.QueryRecordByBondIDRequest{
				Id: suite.bond.GetId(),
			},
			true,
			false,
			1,
		},
	}
	for _, test := range testCasesByBondID {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			resp, err := grpcClient.GetRecordByBondID(context.Background(), test.req)
			if test.expErr {
				sr.Zero(resp.GetRecords())
			} else {
				sr.NoError(err)
				sr.NotNil(resp.GetRecords())
				if test.createRecord {
					sr.NotZero(resp.GetRecords())
					sr.Equal(resp.GetRecords()[0].GetBondId(), suite.bond.GetId())
				}
			}
		})
	}
}

func getAttributesFromRecord(recordAttr map[string]interface{}, key string) interface{} {
	keys := strings.Split(key, "---")
	for i, str := range keys {
		if i == len(keys)-2 {
			return recordAttr[str].(string)
		}
		recordAttr = recordAttr[str].(map[string]interface{})
	}
	return nil
}

func (suite *KeeperTestSuite) TestGrpcQueryRegistryModuleBalance() {
	grpcClient, ctx := suite.queryClient, suite.ctx
	sr := suite.Require()
	examples := []string{
		"/../helpers/examples/service_provider_example.yml",
		"/../helpers/examples/website_registration_example.yml",
	}
	testCases := []struct {
		msg           string
		req           *registrytypes.GetRegistryModuleBalanceRequest
		createRecords bool
		expErr        bool
		noOfRecords   int
	}{
		{
			"Get Module Balance",
			&registrytypes.GetRegistryModuleBalanceRequest{},
			true,
			false,
			1,
		},
	}
	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			if test.createRecords {
				dir, err := os.Getwd()
				sr.NoError(err)
				for _, example := range examples {
					payloadType, err := cli.GetPayloadFromFile(fmt.Sprint(dir, example))
					sr.NoError(err)
					payload, err := payloadType.ToPayload()
					sr.NoError(err)
					record, err := suite.app.RegistryKeeper.ProcessSetRecord(ctx, registrytypes.MsgSetRecord{
						BondId:  suite.bond.GetId(),
						Signer:  suite.accounts[0].String(),
						Payload: payload,
					})
					sr.NoError(err)
					sr.NotNil(record.ID)
				}
			}
			resp, err := grpcClient.GetRegistryModuleBalance(context.Background(), test.req)
			if test.expErr {
				suite.Error(err)
			} else {
				sr.NoError(err)
				sr.Equal(test.noOfRecords, len(resp.GetBalances()))
				if test.createRecords {
					balance := resp.GetBalances()[0]
					sr.Equal(balance.AccountName, registrytypes.RecordRentModuleAccountName)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcQueryWhoIS() {
	grpcClient, ctx := suite.queryClient, suite.ctx
	sr := suite.Require()
	authorityName := "TestGrpcQueryWhoIS"

	testCases := []struct {
		msg         string
		req         *registrytypes.QueryWhoisRequest
		createName  bool
		expErr      bool
		noOfRecords int
	}{
		{
			"Invalid Request without name",
			&registrytypes.QueryWhoisRequest{},
			false,
			true,
			1,
		},
		{
			"Success",
			&registrytypes.QueryWhoisRequest{},
			true,
			false,
			1,
		},
	}
	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			if test.createName {
				err := suite.app.RegistryKeeper.ProcessReserveAuthority(ctx, registrytypes.MsgReserveAuthority{
					Name:   authorityName,
					Signer: suite.accounts[0].String(),
					Owner:  suite.accounts[0].String(),
				})
				sr.NoError(err)
				test.req = &registrytypes.QueryWhoisRequest{Name: authorityName}
			}
			resp, err := grpcClient.Whois(context.Background(), test.req)
			if test.expErr {
				sr.Zero(len(resp.NameAuthority.AuctionId))
			} else {
				sr.NoError(err)
				if test.createName {
					nameAuth := resp.NameAuthority
					sr.NotNil(nameAuth)
					sr.Equal(nameAuth.OwnerAddress, suite.accounts[0].String())
					sr.Equal(registrytypes.AuthorityActive, nameAuth.Status)
				}
			}
		})
	}
}
