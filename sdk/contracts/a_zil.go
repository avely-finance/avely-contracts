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
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type AZil struct {
	Contract
}

func (a *AZil) ChangeAdmin(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"new_admin",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return a.Call("ChangeAdmin", args, "0")
}

func (a *AZil) ClaimAdmin() (*transaction.Transaction, error) {
	args := []core.ContractValue{}
	return a.Call("ClaimAdmin", args, "0")
}

func (b *AZil) ChangeZimplAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return b.Call("ChangeZimplAddress", args, "0")
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

func (a *AZil) ChangeHolderAddress(new_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + new_addr,
		},
	}
	return a.Contract.Call("ChangeHolderAddress", args, "0")
}

func (a *AZil) DelegateStake(amount, initiator string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"initiator",
			"ByStr20",
			"0x" + initiator,
		},
	}

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

func (a *AZil) WithdrawStakeAmt(amount, initiator string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"amount",
			"Uint128",
			amount,
		},
		{
			"initiator",
			"ByStr20",
			"0x" + initiator,
		},
	}
	return a.Call("WithdrawStakeAmt", args, "0")
}

func (a *AZil) DrainBuffer(buffer_addr string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"buffer_addr",
			"ByStr20",
			"0x" + buffer_addr,
		},
	}
	return a.Call("DrainBuffer", args, "0")
}

func (a *AZil) CompleteWithdrawal(initiator string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"initiator",
			"ByStr20",
			"0x" + initiator,
		},
	}
	return a.Call("CompleteWithdrawal", args, "0")
}

func (a *AZil) ZilBalanceOf(addr, initiator string) (*transaction.Transaction, error) {
	args := []core.ContractValue{
		{
			"address",
			"ByStr20",
			"0x" + addr,
		},
		{
			"initiator",
			"ByStr20",
			"0x" + initiator,
		},
	}
	return a.Contract.Call("ZilBalanceOf", args, "0")
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

func NewAZilContract(sdk *AvelySDK, aproxyAddr, zimplAddr string) (*AZil, error) {
	contract := buildAZilContract(sdk, aproxyAddr, zimplAddr)

	tx, err := sdk.DeployTo(&contract)
	if err != nil {
		return nil, err
	}
	tx.Confirm(tx.ID, sdk.Cfg.TxConfrimMaxAttempts, sdk.Cfg.TxConfirmIntervalSec, contract.Provider)
	if tx.Status == core.Confirmed {
		b32, _ := bech32.ToBech32Address(tx.ContractAddress)

		sdkContract := Contract{
			Sdk:             sdk,
			Provider:        *contract.Provider,
			Addr:            tx.ContractAddress,
			Bech32:          b32,
			Wallet:          contract.Signer,
			StateFieldTypes: buildAZilStateFields(),
		}
		return &AZil{Contract: sdkContract}, nil
	} else {
		data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
		return nil, errors.New("deploy failed:" + string(data))
	}
}

func RestoreAZilContract(sdk *AvelySDK, contractAddress, aproxyAddr, zimplAddr string) (*AZil, error) {
	contract := buildAZilContract(sdk, aproxyAddr, zimplAddr)

	b32, err := bech32.ToBech32Address(contractAddress)

	if err != nil {
		return nil, errors.New("Config has invalid AZil address")
	}

	sdkContract := Contract{
		Sdk:             sdk,
		Provider:        *contract.Provider,
		Addr:            contractAddress,
		Bech32:          b32,
		Wallet:          contract.Signer,
		StateFieldTypes: buildAZilStateFields(),
	}
	return &AZil{Contract: sdkContract}, nil
}

func buildAZilContract(sdk *AvelySDK, aproxyAddr, zimplAddr string) contract2.Contract {
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
			Value: "0x" + sdk.GetAddressFromPrivateKey(key),
		}, {
			VName: "init_aproxy_address",
			Type:  "ByStr20",
			Value: "0x" + aproxyAddr,
		}, {
			VName: "init_azil_ssn_address",
			Type:  "ByStr20",
			Value: aZilSSNAddress,
		}, {
			VName: "init_zimpl_address",
			Type:  "ByStr20",
			Value: "0x" + zimplAddr,
		}, {
			VName: "init_holder_address",
			Type:  "ByStr20",
			Value: "0xb2e2c996e6068f4ae11c4cc2c6a189b774819f79",
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

func buildAZilStateFields() StateFieldTypes {
	stateFieldTypes := make(StateFieldTypes)
	stateFieldTypes["balances"] = "StateFieldMap"
	stateFieldTypes["last_buf_deposit_cycle_deleg"] = "StateFieldMap"
	stateFieldTypes["buffers_addresses"] = "StateFieldArray"
	stateFieldTypes["withdrawal_pending"] = "StateFieldMapMapPair"
	stateFieldTypes["withdrawal_unbonded"] = "StateFieldMapPair"
	stateFieldTypes["staging_admin_address"] = "StateFieldOption"

	return stateFieldTypes
}
