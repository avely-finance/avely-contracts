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

type Zimpl struct {
	Contract
}
func NewZimpl(sdk *AvelySDK, ZproxyAddr, GzilAddr string) (*Zimpl, error) {
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
			Value: "0x" + sdk.GetAddressFromPrivateKey(key),
		}, {
			VName: "init_proxy_address",
			Type:  "ByStr20",
			Value: "0x" + ZproxyAddr,
		},
		{
			VName: "init_gzil_address",
			Type:  "ByStr20",
			Value: "0x" + GzilAddr,
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)

	contract := contract2.Contract{
		Code:   string(code),
		Init:   init,
		Signer: wallet,
	}

	tx, err := sdk.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, sdk.Cfg.TxConfrimMaxAttempts, sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		stateFieldTypes := make(StateFieldTypes)
		stateFieldTypes["buff_deposit_deleg"] = "StateFieldMapMapMap"
		stateFieldTypes["direct_deposit_deleg"] = "StateFieldMapMapMap"

		contract := Contract{
			Sdk:             sdk,
			Provider:        *contract.Provider,
			Addr:            tx.ContractAddress,
			Bech32:          b32,
			Wallet:          wallet,
			StateFieldTypes: stateFieldTypes,
		}
		// TxIdLast = tx.ID

		return &Zimpl{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}
