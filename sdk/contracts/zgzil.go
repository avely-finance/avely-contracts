package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	. "github.com/avely-finance/avely-contracts/sdk/core"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
)

type Gzil struct {
	Contract
}

func NewGzil(sdk *AvelySDK, deployer *account.Wallet) (*Gzil, error) {
	code, _ := ioutil.ReadFile("contracts/zilliqa_staking/gzil.scilla")

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		},
		{
			VName: "contract_owner",
			Type:  "ByStr20",
			Value: "0x" + deployer.DefaultAccount.Address,
		},
		{
			VName: "init_minter",
			Type:  "ByStr20",
			Value: "0x" + deployer.DefaultAccount.Address,
		},
		{
			VName: "name",
			Type:  "String",
			Value: "Governance ZIL",
		},
		{
			VName: "symbol",
			Type:  "String",
			Value: "gzil",
		},
		{
			VName: "decimals",
			Type:  "Uint32",
			Value: "15",
		},
		{
			VName: "init_supply",
			Type:  "Uint128",
			Value: "0",
		},
		{
			VName: "num_minting_blocks",
			Type:  "Uint128",
			Value: "0", //was 620500. Minting is over, so we don't need to assume it in tests
		},
	}

	contract := contract2.Contract{
		Provider: sdk.InitProvider(),
		Code:     string(code),
		Init:     init,
		Signer:   deployer,
	}

	tx, err := sdk.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, sdk.Cfg.TxConfrimMaxAttempts, sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		sdkContract := Contract{
			Sdk:      sdk,
			Provider: *contract.Provider,
			Addr:     "0x" + tx.ContractAddress,
			Bech32:   b32,
			Wallet:   deployer,
		}
		sdkContract.ErrorCodes = sdkContract.ParseErrorCodes(contract.Code)
		return &Gzil{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed: " + string(data))
	}
}
