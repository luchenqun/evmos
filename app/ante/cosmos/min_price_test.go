package cosmos_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	cosmosante "github.com/evmos/evmos/v15/app/ante/cosmos"
	"github.com/evmos/evmos/v15/testutil"
	testutiltx "github.com/evmos/evmos/v15/testutil/tx"
	"github.com/evmos/evmos/v15/utils"
)

var execTypes = []struct {
	name      string
	isCheckTx bool
	simulate  bool
}{
	{"checkTx", true, false},
	{"deliverTx", false, false},
	{"deliverTxSimulate", false, true},
}

func (suite *AnteTestSuite) TestMinGasPriceDecorator() {
	denom := utils.BaseDenom
	testMsg := banktypes.MsgSend{
		FromAddress: "evmos1x8fhpj9nmhqk8z9kpgjt95ck2xwyue0ptzkucp",
		ToAddress:   "evmos1dx67l23hz9l0k9hcher8xz04uj7wf3yu26l2yn",
		Amount:      sdk.Coins{sdk.Coin{Amount: sdkmath.NewInt(10), Denom: denom}},
	}

	testCases := []struct {
		name                string
		malleate            func() sdk.Tx
		expPass             bool
		errMsg              string
		allowPassOnSimulate bool
		allowPassOnCheck    bool
	}{
		{
			"invalid cosmos tx type",
			func() sdk.Tx {
				return &testutiltx.InvalidTx{}
			},
			false,
			"invalid transaction type",
			false,
			false,
		},
		{
			"valid cosmos tx with MinGasPrices = 0, gasPrice = 0",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.ZeroDec()
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(0), denom, &testMsg)
				return txBuilder.GetTx()
			},
			true,
			"provided fee < minimum global fee",
			true,
			false,
		},
		{
			"valid cosmos tx with MinGasPrices = 0, gasPrice > 0",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.ZeroDec()
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), denom, &testMsg)
				return txBuilder.GetTx()
			},
			true,
			"",
			true,
			true,
		},
		{
			"valid cosmos tx with MinGasPrices = 10, gasPrice = 10",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.NewDec(10)
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), denom, &testMsg)
				return txBuilder.GetTx()
			},
			true,
			"",
			true,
			true,
		},
		{
			"invalid cosmos tx with MinGasPrices = 10, gasPrice = 0",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.NewDec(10)
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(0), denom, &testMsg)
				return txBuilder.GetTx()
			},
			false,
			"provided fee < minimum global fee",
			true,
			false,
		},
		{
			"invalid cosmos tx with wrong denom",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.NewDec(10)
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(10), "stake", &testMsg)
				return txBuilder.GetTx()
			},
			false,
			"provided fee < minimum global fee",
			true,
			false,
		},
		{
			"valid cosmos tx with MinGasPrices = 0, gasPrice = 0, LocalMinGasPrices = 1",
			func() sdk.Tx {
				params := suite.app.FeeMarketKeeper.GetParams(suite.ctx)
				params.MinGasPrice = sdk.NewDec(0)
				err := suite.app.FeeMarketKeeper.SetParams(suite.ctx, params)
				suite.Require().NoError(err)

				txBuilder := suite.CreateTestCosmosTxBuilder(sdkmath.NewInt(0), denom, &testMsg)
				return txBuilder.GetTx()
			},
			true,
			"",
			true,
			false,
		},
	}

	for _, et := range execTypes {
		for _, tc := range testCases {
			suite.Run(et.name+"_"+tc.name, func() {
				// s.SetupTest(et.isCheckTx)
				ctx := suite.ctx.WithIsReCheckTx(et.isCheckTx)
				dec := cosmosante.NewMinGasPriceDecorator(suite.app.FeeMarketKeeper, suite.app.EvmKeeper)
				_, err := dec.AnteHandle(ctx, tc.malleate(), et.simulate, testutil.NextFn)

				if (tc.expPass && et.name == "deliverTx") ||
					(et.simulate && tc.allowPassOnSimulate && et.name == "deliverTxSimulate") ||
					(tc.allowPassOnCheck && et.name == "checkTx") {
					suite.Require().NoError(err, tc.name)
				} else {
					suite.Require().Error(err, tc.name)
					suite.Require().Contains(err.Error(), tc.errMsg, tc.name)
				}
			})
		}
	}
}
