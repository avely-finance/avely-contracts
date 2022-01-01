package contracts

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"

    . "github.com/avely-finance/avely-contracts/sdk/core"

    "github.com/Zilliqa/gozilliqa-sdk/account"
    "github.com/Zilliqa/gozilliqa-sdk/bech32"
    contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
    "github.com/Zilliqa/gozilliqa-sdk/core"
    provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
    "github.com/Zilliqa/gozilliqa-sdk/transaction"
)

type AZilProxy struct {
    Contract
}

func (a *AZilProxy) CompleteWithdrawal() (*transaction.Transaction, error) {
    args := []core.ContractValue{}
    return a.Call("CompleteWithdrawal", args, "0")
}

func (a *AZilProxy) DelegateStake(amount string) (*transaction.Transaction, error) {
    args := []core.ContractValue{}
    return a.Call("DelegateStake", args, amount)
}

func (a *AZilProxy) UpgradeTo(newImplementation string) (*transaction.Transaction, error) {
    args := []core.ContractValue{
        {
            "newImplementation",
            "ByStr20",
            newImplementation,
        },
    }
    return a.Call("UpgradeTo", args, "0")
}

func (a *AZilProxy) WithdrawStakeAmt(amount string) (*transaction.Transaction, error) {
    args := []core.ContractValue{
        {
            "amount",
            "Uint128",
            amount,
        },
    }
    return a.Call("WithdrawStakeAmt", args, "0")
}

func (a *AZilProxy) ZilBalanceOf(addr string) (string, error) {
    args := []core.ContractValue{
        {
            "address",
            "ByStr20",
            "0x" + addr,
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
            if param.VName == "address" && param.Value != "0x"+addr {
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

func NewAZilProxyContract(sdk *AvelySDK, aimplAddr string) (*AZilProxy, error) {
    contract := buildAZilProxyContract(sdk, aimplAddr)

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
            StateFieldTypes: buildAZilProxyStateFields(),
        }
        return &AZilProxy{Contract: sdkContract}, nil
    } else {
        data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
        return nil, errors.New("deploy failed:" + string(data))
    }
}

func RestoreAZilProxyContract(sdk *AvelySDK, contractAddress, aimplAddr string) (*AZilProxy, error) {
    contract := buildAZilProxyContract(sdk, aimplAddr)

    b32, err := bech32.ToBech32Address(contractAddress)

    if err != nil {
        return nil, errors.New("Config has invalid AZilProxy address")
    }

    sdkContract := Contract{
        Sdk:             sdk,
        Provider:        *contract.Provider,
        Addr:            contractAddress,
        Bech32:          b32,
        Wallet:          contract.Signer,
        StateFieldTypes: buildAZilProxyStateFields(),
    }
    return &AZilProxy{Contract: sdkContract}, nil
}

func buildAZilProxyContract(sdk *AvelySDK, aimplAddr string) contract2.Contract {
    code, _ := ioutil.ReadFile("contracts/aZilProxy.scilla")
    key := sdk.Cfg.AdminKey

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
            VName: "init_aimpl_address",
            Type:  "ByStr20",
            Value: "0x" + aimplAddr,
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

func buildAZilProxyStateFields() StateFieldTypes {
    stateFieldTypes := make(StateFieldTypes)
    return stateFieldTypes
}
