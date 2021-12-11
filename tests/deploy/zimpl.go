package deploy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	//"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type Zimpl struct {
	Contract
}

func NewZimpl(key, ZproxyAddr, GzilAddr string) (*Zimpl, error) {
	code, _ := ioutil.ReadFile("../contracts/zilliqa_staking/ssnlist.scilla")

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin",
			Type:  "ByStr20",
			Value: "0x" + getAddressFromPrivateKey(key),
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

	tx, err := DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, TX_CONFIRM_MAX_ATTEMPTS, TX_CONFIRM_INTERVAL_SEC, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		stateFieldTypes := make(StateFieldTypes)
		stateFieldTypes["buff_deposit_deleg"] = "StateFieldMapMapMap"
		stateFieldTypes["direct_deposit_deleg"] = "StateFieldMapMapMap"

		contract := Contract{
			Code:            string(code),
			Init:            init,
			Addr:            tx.ContractAddress,
			Bech32:          b32,
			Wallet:          wallet,
			StateFieldTypes: stateFieldTypes,
		}
		TxIdLast = tx.ID

		return &Zimpl{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		log.Println(string(data))
		return nil, errors.New("deploy failed")
	}
}
