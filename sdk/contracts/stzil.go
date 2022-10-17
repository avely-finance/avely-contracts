package contracts

import (
	"encoding/json"
	"errors"
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

type StZIL struct {
	Contract
}

func (a *StZIL) WithUser(key string) *StZIL {
	wallet := account.NewWallet()
	wallet.AddByPrivateKey(key)
	a.Contract.Wallet = wallet
	return a
}

func (s *StZIL) BalanceOf(addr string) *big.Int {
	rawState := s.Contract.SubState("balances", []string{addr})
	state := NewState(rawState)

	return state.Dig("result.balances." + addr).BigInt()
}

func (s *StZIL) IncreaseAllowance(spender, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "spender",
			Type:  "ByStr20",
			Value: spender,
		}, {
			VName: "amount",
			Type:  "Uint128",
			Value: amount,
		},
	}

	return s.Call("IncreaseAllowance", args, "0")
}

func (s *StZIL) DecreaseAllowance(spender, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "spender",
			Type:  "ByStr20",
			Value: spender,
		}, {
			VName: "amount",
			Type:  "Uint128",
			Value: amount,
		},
	}

	return s.Call("DecreaseAllowance", args, "0")
}

func (s *StZIL) Transfer(to, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "to",
			Type:  "ByStr20",
			Value: to,
		}, {
			VName: "amount",
			Type:  "Uint128",
			Value: amount,
		},
	}

	return s.Call("Transfer", args, "0")
}

func (s *StZIL) TransferFrom(from, to, amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			VName: "from",
			Type:  "ByStr20",
			Value: from,
		}, {
			VName: "to",
			Type:  "ByStr20",
			Value: to,
		}, {
			VName: "amount",
			Type:  "Uint128",
			Value: amount,
		},
	}

	return s.Call("TransferFrom", args, "0")
}

func (a *StZIL) ChangeAdmin(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_admin",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeAdmin", args, "0")
}

func (a *StZIL) ClaimAdmin() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("ClaimAdmin", args, "0")
}

func (a *StZIL) ChangeZimplAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeZimplAddress", args, "0")
}

func (a *StZIL) ChangeOwner(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_owner",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeOwner", args, "0")
}

func (a *StZIL) ClaimOwner() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("ClaimOwner", args, "0")
}

// returns
// {"id":"1","jsonrpc":"2.0","result":{
//			"buffer_drained_cycle":
//	  						 {"0x79c7e38dd3b3c88a3fb182f26b66d8889e61cbd6":"123",
//                  "0xbfb3bbde860bcd17315ec0e171ac971de7bea9a3":"124"}
// }
func (a *StZIL) GetDrainedBuffers() map[string]gjson.Result {
	rawState := a.Contract.SubState("buffer_drained_cycle", []string{})
	state := NewState(rawState)
	return state.Dig("result.buffer_drained_cycle").Map()
}

func (a *StZIL) GetAutorestakeAmount() *big.Int {
	rawState := a.Contract.SubState("autorestakeamount", []string{})
	state := NewState(rawState)

	return state.Dig("result.autorestakeamount").BigInt()
}

func (a *StZIL) GetSsnWhitelist() []string {
	partialState := a.Contract.SubState("ssn_addresses", []string{})
	state := NewState(partialState)
	ssnAddrs := state.Dig("result.ssn_addresses").ArrayString()
	return ssnAddrs
}

func (a *StZIL) GetSsnIndex() *big.Int {
	rawState := a.Contract.SubState("ssn_index", []string{})
	state := NewState(rawState)
	return state.Dig("result.ssn_index").BigInt()
}

func (a *StZIL) GetStZilPrice() *big.Float {
	params := a.Contract.BuildBatchParams([]string{"total_supply", "totalstakeamount"})
	raw, _ := a.Contract.BatchSubState(params)
	state := NewState(raw)

	total_supply := state.Dig("0.result.total_supply").BigFloat()
	if 0 == big.NewFloat(0).Cmp(total_supply) {
		return big.NewFloat(-1)
	}
	totalstakeamount := state.Dig("1.result.totalstakeamount").BigFloat()

	return DivBF(totalstakeamount, total_supply)
}

func (s *StZIL) GetTreasuryAddress() string {
	rawState := s.Contract.SubState("treasury_address", []string{})
	state := NewState(rawState)
	return state.Dig("result.treasury_address").String()
}

func (a *StZIL) ChangeBuffers(newBuffers []string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_buffers",
			"List ByStr20",
			newBuffers,
		},
	}
	return a.Contract.Call("ChangeBuffers", args, "0")
}

func (a *StZIL) AddSSN(ssnaddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
	}
	return a.Contract.Call("AddSSN", args, "0")
}

func (a *StZIL) RemoveSSN(ssnaddr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"ssnaddr",
			"ByStr20",
			ssnaddr,
		},
	}
	return a.Contract.Call("RemoveSSN", args, "0")
}

func (a *StZIL) ClaimWithdrawal(ready_blocks []string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"blocks_to_withdraw",
			"List BNum",
			ready_blocks,
		},
	}
	return a.Contract.Call("ClaimWithdrawal", args, "0")
}

func (a *StZIL) ChangeTreasuryAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			new_addr,
		},
	}
	return a.Call("ChangeTreasuryAddress", args, "0")
}

func (a *StZIL) SetHolderAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			new_addr,
		},
	}
	return a.Contract.Call("SetHolderAddress", args, "0")
}

func (a *StZIL) ChangeRewardsFee(new_fee string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_fee",
			"Uint128",
			new_fee,
		},
	}
	return a.Call("ChangeRewardsFee", args, "0")
}

func (a *StZIL) ChownStakeConfirmSwap(delegator string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"delegator",
			"ByStr20",
			delegator,
		},
	}
	return a.Call("ChownStakeConfirmSwap", args, "0")
}

func (a *StZIL) ChownStakeReDelegate(from_ssn, amount string) (*transaction.Transaction, error) {
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

func (a *StZIL) DelegateStake(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("DelegateStake", args, amount)
}

func (a *StZIL) IncreaseAutoRestakeAmount(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("IncreaseAutoRestakeAmount", args, amount)
}

func (a *StZIL) PerformAutoRestake() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("PerformAutoRestake", args, "0")
}

func (a *StZIL) UpdateStakingParameters(min_deleg_stake string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"min_deleg_stake",
			"Uint128",
			min_deleg_stake,
		},
	}
	return a.Call("UpdateStakingParameters", args, "0")
}

func (a *StZIL) WithdrawStakeAmt(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return a.Call("WithdrawStakeAmt", args, "0")
}

func (a *StZIL) CompleteWithdrawal() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("CompleteWithdrawal", args, "0")
}

func (a *StZIL) ConsolidateInHolder(buffer_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"buffer_addr",
			"ByStr20",
			buffer_addr,
		},
	}
	return a.Call("ConsolidateInHolder", args, "0")
}

func (a *StZIL) ZilBalanceOf(addr string) *big.Int {
	stzilPriceFloat := a.GetStZilPrice()
	balance := a.BalanceOf(addr)
	balanceFloat := new(big.Float).SetInt(balance)
	zilBalanceFloat := new(big.Float).Mul(stzilPriceFloat, balanceFloat)

	result := new(big.Int)
	zilBalanceFloat.Int(result) // store converted number in result

	return result
}

func (a *StZIL) ClaimRewards(address, ssn string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"buffer_or_holder",
			"ByStr20",
			address,
		},
		{
			"ssn",
			"ByStr20",
			ssn,
		},
	}
	return a.Call("ClaimRewards", args, "0")
}

func (a *StZIL) ClaimRewardsSuccessCallBack() (*transaction.Transaction, error) {
	args := []core.ContractValue{}

	return a.Call("ClaimRewardsSuccessCallBack", args, "0")
}

func (a *StZIL) DelegateStakeSuccessCallBack(amount string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
	}
	return a.Call("DelegateStakeSuccessCallBack", args, "0")
}

func (a *StZIL) CompleteWithdrawalSuccessCallBack() (*transaction.Transaction, error) {
	args := []core.ContractValue{}

	return a.Call("CompleteWithdrawalSuccessCallBack", args, "0")
}

func (a *StZIL) PauseIn() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("PauseIn", args, "0")
}

func (a *StZIL) UnpauseIn() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("UnPauseIn", args, "0")
}

func (a *StZIL) PauseOut() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("PauseOut", args, "0")
}

func (a *StZIL) UnpauseOut() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("UnPauseOut", args, "0")
}

func (a *StZIL) PauseZrc2() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("PauseZrc2", args, "0")
}

func (a *StZIL) UnpauseZrc2() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("UnPauseZrc2", args, "0")
}

func NewStZILContract(sdk *AvelySDK, owner, zimplAddr string) (*StZIL, error) {
	contract := buildStZILContract(sdk, owner, zimplAddr)

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
		return &StZIL{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreStZILContract(sdk *AvelySDK, contractAddress, owner, zimplAddr string) (*StZIL, error) {
	contract := buildStZILContract(sdk, owner, zimplAddr)

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid StZIL address")
	}

	sdkContract := Contract{
		Sdk:      sdk,
		Provider: *contract.Provider,
		Addr:     contractAddress,
		Bech32:   b32,
		Wallet:   contract.Signer,
	}
	return &StZIL{Contract: sdkContract}, nil
}

func buildStZILContract(sdk *AvelySDK, owner, zimplAddr string) contract2.Contract {
	code, _ := ioutil.ReadFile("contracts/stzil.scilla")

	init := []core.ContractValue{
		{
			VName: "_scilla_version",
			Type:  "Uint32",
			Value: "0",
		}, {
			VName: "contract_owner",
			Type:  "ByStr20",
			Value: owner,
		}, {
			VName: "init_admin_address",
			Type:  "ByStr20",
			Value: sdk.GetAddressFromPrivateKey(sdk.Cfg.AdminKey),
		}, {
			VName: "init_zimpl_address",
			Type:  "ByStr20",
			Value: zimplAddr,
		}, {
			VName: "name",
			Type:  "String",
			Value: "StZIL",
		}, {
			VName: "symbol",
			Type:  "String",
			Value: "stZIL",
		}, {
			VName: "decimals",
			Type:  "Uint32",
			Value: "12",
		}, {
			VName: "init_supply",
			Type:  "Uint128",
			Value: "0",
		},
	}

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(sdk.Cfg.AdminKey)

	return contract2.Contract{
		Provider: sdk.InitProvider(),
		Code:     string(code),
		Init:     init,
		Signer:   wallet,
	}
}
