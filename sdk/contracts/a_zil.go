package contracts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/tidwall/gjson"

	. "github.com/avely-finance/avely-contracts/sdk/core"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type AZil struct {
	Contract
}

func (a *AZil) WithUser(key string) *AZil {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)
	a.Contract.Wallet = wallet
	return a
}

func (a *AZil) ChangeAdmin(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_admin",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeAdmin", args, "0")
}

func (a *AZil) ClaimAdmin() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("ClaimAdmin", args, "0")
}

func (a *AZil) ChangeZimplAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeZimplAddress", args, "0")
}

// returns
// {"id":"1","jsonrpc":"2.0","result":{
//			"buffer_drained_cycle":
//	  						 {"0x79c7e38dd3b3c88a3fb182f26b66d8889e61cbd6":"123",
//                  "0xbfb3bbde860bcd17315ec0e171ac971de7bea9a3":"124"}
// }
func (a *AZil) GetDrainedBuffers() map[string]gjson.Result {
	rawState := a.Contract.SubState("buffer_drained_cycle", []string{})
	state := NewState(rawState)
	return state.Dig("result.buffer_drained_cycle").Map()
}

func (a *AZil) GetAutorestakeAmount() *big.Int {
	rawState := a.Contract.SubState("autorestakeamount", []string{})
	state := NewState(rawState)

	return state.Dig("result.autorestakeamount").BigInt()
}

func (a *AZil) GetAzilPrice() *big.Float {
	params := a.Contract.BuildBatchParams([]string{"totaltokenamount", "totalstakeamount"})
	raw, _ := a.Contract.BatchSubState(params)
	state := NewState(raw)

	totaltokenamount := state.Dig("0.result.totaltokenamount").BigFloat()
	totalstakeamount := state.Dig("1.result.totalstakeamount").BigFloat()

	return DivBF(totalstakeamount, totaltokenamount)
}

func (a *AZil) ChangeBuffers(new_buffers []string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_buffers",
			"List ByStr20",
			new_buffers,
		},
	}
	return a.Contract.Call("ChangeBuffers", args, "0")
}

func (a *AZil) ClaimWithdrawal(ready_blocks []string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"blocks_to_withdraw",
			"List BNum",
			ready_blocks,
		},
	}
	return a.Contract.Call("ClaimWithdrawal", args, "0")
}

func (a *AZil) ChangeAzilSSNAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeAzilSSNAddress", args, "0")
}

func (a *AZil) ChangeHolderAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			new_addr,
		},
	}
	return a.Contract.Call("ChangeHolderAddress", args, "0")
}

func (a *AZil) ChangeRewardsFee(new_fee string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_fee",
			"Uint128",
			new_fee,
		},
	}
	return a.Call("ChangeRewardsFee", args, "0")
}

func (a *AZil) ChownStakeConfirmSwap(delegator string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"delegator",
			"ByStr20",
			delegator,
		},
	}
	return a.Call("ChownStakeConfirmSwap", args, "0")
}

func (a *AZil) ChownStakeReDelegate(from_ssn, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"from_ssn",
			"ByStr20",
			from_ssn,
		},
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return a.Call("ChownStakeReDelegate", args, "0")
}

func (a *AZil) DelegateStake(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("DelegateStake", args, amount)
}

func (a *AZil) IncreaseAutoRestakeAmount(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("IncreaseAutoRestakeAmount", args, amount)
}

func (a *AZil) PerformAutoRestake() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("PerformAutoRestake", args, "0")
}

func (a *AZil) UpdateStakingParameters(min_deleg_stake string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"min_deleg_stake",
			"Uint128",
			min_deleg_stake,
		},
	}
	return a.Call("UpdateStakingParameters", args, "0")
}

func (a *AZil) WithdrawStakeAmt(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return a.Call("WithdrawStakeAmt", args, "0")
}

func (a *AZil) DrainBuffer(buffer_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"buffer_addr",
			"ByStr20",
			buffer_addr,
		},
	}
	return a.Call("DrainBuffer", args, "0")
}

func (a *AZil) CompleteWithdrawal() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("CompleteWithdrawal", args, "0")
}

func (a *AZil) ZilBalanceOf(addr string) (string, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			addr,
		},
	}
	tx, err := a.Contract.Call("ZilBalanceOf", args, "0")
	if err != nil {
		return "", err
	}

	for _, transition := range tx.Receipt.Transitions {
		if "ZilBalanceOfCallBack" != transition.Msg.Tag {
			continue
		}
		for _, param := range transition.Msg.Params {
			if param.VName == "address" && param.Value != addr {
				//it's balance of some other address, it should not be so
				return "", errors.New("Balance not found for addr=" + addr)
			}
			if param.VName == "balance" {
				return fmt.Sprintf("%v", param.Value), nil
			}
		}
		break
	}
	return "", errors.New("Balance not found")
}

func (a *AZil) ClaimRewardsSuccessCallBack() (*transaction.Transaction, error) {
	args := []core.ContractValue{}

	return a.Call("ClaimRewardsSuccessCallBack", args, "0")
}

func (a *AZil) DelegateStakeSuccessCallBack(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return a.Call("DelegateStakeSuccessCallBack", args, "0")
}

func (a *AZil) CompleteWithdrawalSuccessCallBack() (*transaction.Transaction, error) {
	args := []core.ContractValue{}

	return a.Call("CompleteWithdrawalSuccessCallBack", args, "0")
}

func (a *AZil) Pause() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("Pause", args, "0")
}

func (a *AZil) Unpause() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("UnPause", args, "0")
}

func NewAZilContract(sdk *AvelySDK, zimplAddr string) (*AZil, error) {
	contract := buildAZilContract(sdk, zimplAddr)

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
		return &AZil{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreAZilContract(sdk *AvelySDK, contractAddress, zimplAddr string) (*AZil, error) {
	contract := buildAZilContract(sdk, zimplAddr)

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid AZil address")
	}

	sdkContract := Contract{
		Sdk:      sdk,
		Provider: *contract.Provider,
		Addr:     contractAddress,
		Bech32:   b32,
		Wallet:   contract.Signer,
	}
	return &AZil{Contract: sdkContract}, nil
}

func buildAZilContract(sdk *AvelySDK, zimplAddr string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/aZil.scilla")
	key := sdk.Cfg.AdminKey
	aZilSSNAddress := sdk.Cfg.AzilSsnAddress

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "init_admin_address",
			Type:  "ByStr20",
			Value: sdk.GetAddressFromPrivateKey(key),
		}, {
			VName: "init_azil_ssn_address",
			Type:  "ByStr20",
			Value: aZilSSNAddress,
		}, {
			VName: "init_zimpl_address",
			Type:  "ByStr20",
			Value: zimplAddr,
		}, {
			VName: "init_holder_address",
			Type:  "ByStr20",
			Value: "0xb2e2c996e6068f4ae11c4cc2c6a189b774819f79",
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
