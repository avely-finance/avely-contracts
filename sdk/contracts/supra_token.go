package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

type SupraToken struct {
	Contract
}

func NewSupraToken(sdk *AvelySDK) (*SupraToken, error) {
	contract := buildSwapTokenContract(sdk)

	tx, err := sdk.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, sdk.Cfg.TxConfrimMaxAttempts, sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		contract := Contract{
			Sdk:      sdk,
			Provider: *contract.Provider,
			Addr:     "0x" + tx.ContractAddress,
			Bech32:   b32,
			Wallet:   contract.Signer,
		}

		return &SupraToken{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func buildSwapTokenContract(sdk *AvelySDK) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/zilswap/supra_token.scilla")
	key := sdk.Cfg.AdminKey

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "contract_owner",
			Type:  "ByStr20",
			Value: sdk.GetAddressFromPrivateKey(key),
		}, {
			VName: "name",
			Type:  "String",
			Value: "SUPRA",
		}, {
			VName: "symbol",
			Type:  "String",
			Value: "SUPRA",
		}, {
			VName: "decimals",
			Type:  "Uint32",
			Value: "9",
		}, {
			VName: "init_supply",
			Type:  "Uint128",
			Value: "2000000000000000",
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	return contract2.Contract{
		Provider: sdk.InitProvider(),
		Code:     string(code),
		Init:     init,
		Signer:   wallet,
	}
}
