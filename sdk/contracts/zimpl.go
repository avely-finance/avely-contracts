package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

type Zimpl struct {
	Contract
}

func NewZimpl(sdk *AvelySDK, ZproxyAddr, GzilAddr string) (*Zimpl, error) {
	contract := buildZimplContract(sdk, ZproxyAddr, GzilAddr)

	tx, err := sdk.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, sdk.Cfg.TxConfrimMaxAttempts, sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		contract := Contract{
			Sdk:             sdk,
			Provider:        *contract.Provider,
			Addr:            "0x" + tx.ContractAddress,
			Bech32:          b32,
			Wallet:          contract.Signer,
			StateFieldTypes: buildZimplStateFields(),
		}

		return &Zimpl{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreZimpl(sdk *AvelySDK, contractAddress, ZproxyAddr, GzilAddr string) (*Zimpl, error) {
	contract := buildZimplContract(sdk, ZproxyAddr, GzilAddr)

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid Zimpl address")
	}

	sdkContract := Contract{
		Sdk:             sdk,
		Provider:        *contract.Provider,
		Addr:            contractAddress,
		Bech32:          b32,
		Wallet:          contract.Signer,
		StateFieldTypes: buildZimplStateFields(),
	}

	return &Zimpl{Contract: sdkContract}, nil
}

func buildZimplContract(sdk *AvelySDK, ZproxyAddr, GzilAddr string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/zilliqa_staking/ssnlist.scilla")
	key := sdk.Cfg.AdminKey

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin",
			Type:  "ByStr20",
			Value: sdk.GetAddressFromPrivateKey(key),
		}, {
			VName: "init_proxy_address",
			Type:  "ByStr20",
			Value: ZproxyAddr,
		},
		{
			VName: "init_gzil_address",
			Type:  "ByStr20",
			Value: GzilAddr,
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	return contract2.Contract{
		Provider: provider2.NewProvider(sdk.Cfg.ApiUrl),
		Code:     string(code),
		Init:     init,
		Signer:   wallet,
	}
}

func buildZimplStateFields() StateFieldTypes {
	stateFieldTypes := make(StateFieldTypes)
	stateFieldTypes["buff_deposit_deleg"] = "StateFieldMapMapMap"
	stateFieldTypes["direct_deposit_deleg"] = "StateFieldMapMapMap"

	return stateFieldTypes
}
