package keeper_test

import (
	"context"
	"fmt"
	"github.com/tharsis/ethermint/x/nameservice/client/cli"
	nameservicetypes "github.com/tharsis/ethermint/x/nameservice/types"
	"os"
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
	testCases := []struct {
		msg          string
		req          *nameservicetypes.QueryListRecordsRequest
		createRecord bool
		expErr       bool
		noOfRecords  int
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
			1,
		},
	}
	for _, test := range testCases {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			if test.createRecord {
				dir, err := os.Getwd()
				sr.NoError(err)
				payload, err := cli.GetPayloadFromFile(dir + "/../helpers/examples/example1.yml")
				sr.NoError(err)
				err = suite.app.NameServiceKeeper.ProcessSetRecord(ctx, nameservicetypes.MsgSetRecord{
					BondId:  suite.bond.GetId(),
					Signer:  suite.accounts[0].String(),
					Payload: payload.ToPayload(),
				})
				sr.NoError(err)
			}
			resp, err := grpcClient.ListRecords(context.Background(), test.req)
			if test.expErr {
				suite.Error(err)
			} else {
				sr.NoError(err)
				sr.Equal(test.noOfRecords, len(resp.GetRecords()))
				if test.createRecord {
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
		req          *nameservicetypes.QueryRecordByIdRequest
		createRecord bool
		expErr       bool
		noOfRecords  int
	}{
		{
			"Invalid Request without record id",
			&nameservicetypes.QueryRecordByIdRequest{},
			false,
			true,
			0,
		},
		{
			"With Record ID",
			&nameservicetypes.QueryRecordByIdRequest{
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
		req          *nameservicetypes.QueryRecordByBondIdRequest
		createRecord bool
		expErr       bool
		noOfRecords  int
	}{
		{
			"Invalid Request without bond id",
			&nameservicetypes.QueryRecordByBondIdRequest{},
			false,
			true,
			0,
		},
		{
			"With Bond ID",
			&nameservicetypes.QueryRecordByBondIdRequest{
				Id: suite.bond.GetId(),
			},
			true,
			false,
			1,
		},
	}
	for _, test := range testCasesByBondID {
		suite.Run(fmt.Sprintf("Case %s ", test.msg), func() {
			resp, err := grpcClient.GetRecordByBondId(context.Background(), test.req)
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
	testCases := []struct {
		msg          string
		req          *nameservicetypes.GetNameServiceModuleBalanceRequest
		createRecord bool
		expErr       bool
		noOfRecords  int
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
			if test.createRecord {
				dir, err := os.Getwd()
				sr.NoError(err)
				payload, err := cli.GetPayloadFromFile(dir + "/../helpers/examples/example1.yml")
				sr.NoError(err)
				err = suite.app.NameServiceKeeper.ProcessSetRecord(ctx, nameservicetypes.MsgSetRecord{
					BondId:  suite.bond.GetId(),
					Signer:  suite.accounts[0].String(),
					Payload: payload.ToPayload(),
				})
				sr.NoError(err)
			}
			resp, err := grpcClient.GetNameServiceModuleBalance(context.Background(), test.req)
			if test.expErr {
				suite.Error(err)
			} else {
				sr.NoError(err)
				sr.Equal(test.noOfRecords, len(resp.GetBalances()))
				if test.createRecord {
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
	var authorityName = "TestGrpcQueryWhoIS"

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
