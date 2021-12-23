package contracts

import (
	"Azil/test/helpers"

	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
)

type Gzil struct {
	Contract
}

func NewGzil(key string) (*Gzil, error) {
	code, _ := ioutil.ReadFile("../contracts/zilliqa_staking/gzil.scilla")

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		},
		{
			VName: "contract_owner",
			Type:  "ByStr20",
			Value: "0x" + helpers.GetAddressFromPrivateKey(key),
		},
		{
			VName: "init_minter",
			Type:  "ByStr20",
			Value: "0x" + helpers.GetAddressFromPrivateKey(key),
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

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	contract := contract2.Contract{
		Code:   string(code),
		Init:   init,
		Signer: wallet,
	}

	tx, err := helpers.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, TX_CONFIRM_MAX_ATTEMPTS, TX_CONFIRM_INTERVAL_SEC, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		stateFieldTypes := make(StateFieldTypes)

		contract := Contract{
			Provider:        *contract.Provider,
			Addr:            tx.ContractAddress,
			Bech32:          b32,
			Wallet:          wallet,
			StateFieldTypes: stateFieldTypes,
		}
		TxIdLast = tx.ID

		return &Gzil{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed: " + string(data))
	}
}
