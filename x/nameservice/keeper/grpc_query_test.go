package keeper_test

import (
	"context"
	"fmt"
	"os"

	"github.com/cerc-io/laconicd/x/nameservice/client/cli"
	nameservicetypes "github.com/cerc-io/laconicd/x/nameservice/types"
)

func (suite *KeeperTestSuite) TestGrpcQueryParams() {
	grpcClient := suite.queryClient

	testCases := []struct {
		msg string
		req *nameservicetypes.QueryParamsRequest
	}{
		{
			"Get Params",
			&nameservicetypes.QueryParamsRequest{},
		},
	}
	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			resp, _ := grpcClient.Params(context.Background(), test.req)
			defaultParams := nameservicetypes.DefaultParams()
			suite.Require().Equal(defaultParams.String(), resp.GetParams().String())
		})
	}
}

func (suite *KeeperTestSuite) TestGrpcGetRecordLists() {
	grpcClient, ctx := suite.queryClient, suite.ctx
	sr := suite.Require()
	var recordId string
	examples := []string{
		"/../helpers/examples/service_provider_example.yml",
		"/../helpers/examples/website_registration_example.yml",
	}
	testCases := []struct {
		msg           string
		req           *nameservicetypes.QueryListRecordsRequest
		createRecords bool
		expErr        bool
		noOfRecords   int
	}{
		{
			"Empty Records",
			&nameservicetypes.QueryListRecordsRequest{},
			false,
			false,
			0,
		},
		{
			"List Records",
			&nameservicetypes.QueryListRecordsRequest{},
			true,
			false,
			2,
		},
		{
			"Filter with type",
			&nameservicetypes.QueryListRecordsRequest{
				Attributes: []*nameservicetypes.QueryListRecordsRequest_KeyValueInput{
					{
						Key: "type",
						Value: &nameservicetypes.QueryListRecordsRequest_ValueInput{
							Type:    "string",
							String_: "WebsiteRegistrationRecord",
						},
					},
				},
			},
			true,
			false,
			1,
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
					record, err := suite.app.NameServiceKeeper.ProcessSetRecord(ctx, nameservicetypes.MsgSetRecord{
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
				}
			}
		})
	}

	// Get the records by record id
	testCases1 := []struct {
		msg          string
		req          *nameservicetypes.QueryRecordByIDRequest
		createRecord bool
		expErr       bool
		noOfRecords  int
	}{
		{
			"Invalid Request without record id",
			&nameservicetypes.QueryRecordByIDRequest{},
			false,
			true,
			0,
		},
		{
			"With Record ID",
			&nameservicetypes.QueryRecordByIDRequest{
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
		req          *nameservicetypes.QueryRecordByBondIDRequest
		createRecord bool
		expErr       bool
		noOfRecords  int
	}{
		{
			"Invalid Request without bond id",
			&nameservicetypes.QueryRecordByBondIDRequest{},
			false,
			true,
			0,
		},
		{
			"With Bond ID",
			&nameservicetypes.QueryRecordByBondIDRequest{
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

func (suite *KeeperTestSuite) TestGrpcQueryNameserviceModuleBalance() {
	grpcClient, ctx := suite.queryClient, suite.ctx
	sr := suite.Require()
	examples := []string{
		"/../helpers/examples/service_provider_example.yml",
		"/../helpers/examples/website_registration_example.yml",
	}
	testCases := []struct {
		msg           string
		req           *nameservicetypes.GetNameServiceModuleBalanceRequest
		createRecords bool
		expErr        bool
		noOfRecords   int
	}{
		{
			"Get Module Balance",
			&nameservicetypes.GetNameServiceModuleBalanceRequest{},
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
					record, err := suite.app.NameServiceKeeper.ProcessSetRecord(ctx, nameservicetypes.MsgSetRecord{
						BondId:  suite.bond.GetId(),
						Signer:  suite.accounts[0].String(),
						Payload: payload,
					})
					sr.NoError(err)
					sr.NotNil(record.ID)
				}
			}
			resp, err := grpcClient.GetNameServiceModuleBalance(context.Background(), test.req)
			if test.expErr {
				suite.Error(err)
			} else {
				sr.NoError(err)
				sr.Equal(test.noOfRecords, len(resp.GetBalances()))
				if test.createRecords {
					balance := resp.GetBalances()[0]
					sr.Equal(balance.AccountName, nameservicetypes.RecordRentModuleAccountName)
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
		req         *nameservicetypes.QueryWhoisRequest
		createName  bool
		expErr      bool
		noOfRecords int
	}{
		{
			"Invalid Request without name",
			&nameservicetypes.QueryWhoisRequest{},
			false,
			true,
			1,
		},
		{
			"Success",
			&nameservicetypes.QueryWhoisRequest{},
			true,
			false,
			1,
		},
	}
	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			if test.createName {
				err := suite.app.NameServiceKeeper.ProcessReserveAuthority(ctx, nameservicetypes.MsgReserveAuthority{
					Name:   authorityName,
					Signer: suite.accounts[0].String(),
					Owner:  suite.accounts[0].String(),
				})
				sr.NoError(err)
				test.req = &nameservicetypes.QueryWhoisRequest{Name: authorityName}
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
					sr.Equal(nameservicetypes.AuthorityActive, nameAuth.Status)
				}
			}
		})
	}
}
