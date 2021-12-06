package deploy

import (
    "encoding/json"
    "errors"
    //"fmt"
    "io/ioutil"
    "log"

    "github.com/Zilliqa/gozilliqa-sdk/account"
    "github.com/Zilliqa/gozilliqa-sdk/bech32"
    contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
    "github.com/Zilliqa/gozilliqa-sdk/core"
)

type AzilUtils struct {
    Contract
}

func NewAzilUtilsContract(key string) (*AzilUtils, error) {
    code, _ := ioutil.ReadFile("../contracts/azilutils.scillib")
    type Constructor struct {
        Constructor string   `json:"constructor"`
        ArgTypes    []string `json:"argtypes"`
        Arguments   []string `json:"arguments"`
    }
    argtypes := make([]string, 0)
    arguments := make([]string, 0)
    cons := Constructor{
        Constructor: "True",
        ArgTypes:    argtypes,
        Arguments:   arguments,
    }
    init := []core.ContractValue{
        {
            VName: "_scilla_version",
            Type:  "Uint32",
            Value: "0",
        }, {
            VName: "_library",
            Type:  "Bool",
            Value: cons,
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
    tx.Confirm(tx.ID, TxConfirmMaxAttempts, TxConfirmInterval, contract.Provider)
    if tx.Status == core.Confirmed {
        b32, _ := bech32.ToBech32Address(tx.ContractAddress)
        contract := Contract{
            Code:   string(code),
            Init:   init,
            Addr:   tx.ContractAddress,
            Bech32: b32,
            Wallet: wallet,
        }

        return &AzilUtils{Contract: contract}, nil
    } else {
        data, _ := json.MarshalIndent(tx.Receipt, "", "     ")
        log.Println(string(data))
        return nil, errors.New("deploy failed")
    }
}
