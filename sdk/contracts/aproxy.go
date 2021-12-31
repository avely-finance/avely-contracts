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
)

type AZilProxy struct {
    Contract
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
