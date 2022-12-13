package contracts

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

type Zimpl struct {
	Contract
}

func (z *Zimpl) GetSsnList() []string {
	partialState := z.Contract.SubState("ssnlist", []string{})
	state := NewState(partialState)
	ssnAddrs := state.Dig("result.ssnlist|@keys").ArrayString()
	return ssnAddrs
}

func (z *Zimpl) GetDepositAmtDeleg(delegator string) map[string]*big.Int {
	delegator = strings.ToLower(delegator)
	rawState := z.Contract.SubState("deposit_amt_deleg", []string{delegator})
	state := NewState(rawState)
	stateItem := state.Dig("result.deposit_amt_deleg." + delegator)
	return stateItem.MapAddressAmount()
}

func (z *Zimpl) GetBnumReq() int {
	partialState := z.Contract.SubState("bnum_req", []string{})
	state := NewState(partialState)
	bnumReq := state.Dig("result.bnum_req").Int()
	return int(bnumReq)
}

func (z *Zimpl) GetLastRewardCycle() int {
	partialState := z.Contract.SubState("lastrewardcycle", []string{})

	state := NewState(partialState)

	lrc := state.Dig("result.lastrewardcycle").Int()

	return int(lrc)
}

func NewZimpl(sdk *AvelySDK, ZproxyAddr, GzilAddr string, deployer *account.Wallet) (*Zimpl, error) {
	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin",
			Type:  "ByStr20",
			Value: "0x" + deployer.DefaultAccount.Address,
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

	contract := buildZimplContract(sdk, init)
	contract.Signer = deployer

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
			Wallet:   contract.Signer,
		}
		sdkContract.ErrorCodes = sdkContract.ParseErrorCodes(contract.Code)

		return &Zimpl{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreZimpl(sdk *AvelySDK, contractAddress string) (*Zimpl, error) {
	contract := buildZimplContract(sdk, []core.ContractValue{})

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid Zimpl address")
	}

	sdkContract := Contract{
		Sdk:      sdk,
		Provider: *contract.Provider,
		Addr:     contractAddress,
		Bech32:   b32,
		Wallet:   contract.Signer,
	}
	sdkContract.ErrorCodes = sdkContract.ParseErrorCodes(contract.Code)

	return &Zimpl{Contract: sdkContract}, nil
}

func buildZimplContract(sdk *AvelySDK, init []core.ContractValue) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/zilliqa_staking/ssnlist.scilla")

	return contract2.Contract{
		Provider: sdk.InitProvider(),
		Code:     string(code),
		Init:     init,
		Signer:   nil,
	}
}
