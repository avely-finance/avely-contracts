package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

type MinterProxy struct {
	Contract
}

func (m *MinterProxy) WithUser(key string) *MinterProxy {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)
	m.Contract.Wallet = wallet
	return m
}

func (m *MinterProxy) Mint(zilAmount string, minAzilAmount string, deadlineBlock string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "min_azil_amount",
			Type:  "Uint128",
			Value: minAzilAmount,
		}, {
			VName: "deadline_block",
			Type:  "BNum",
			Value: deadlineBlock,
		},
	}

	return m.Call("Mint", args, zilAmount)
}

func NewMinterProxy(sdk *AvelySDK, azilAddr, zilSwapAddr string) (*MinterProxy, error) {
	contract := buildMinterProxyContract(sdk, azilAddr, zilSwapAddr)

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

		return &MinterProxy{Contract: contract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func buildMinterProxyContract(sdk *AvelySDK, azilAddr, zilSwapAddr string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/minter_proxy.scilla")
	key := sdk.Cfg.AdminKey

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_azil_address",
			Type:  "ByStr20",
			Value: azilAddr,
		}, {
			VName: "init_zilswap_address",
			Type:  "ByStr20",
			Value: zilSwapAddr,
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
