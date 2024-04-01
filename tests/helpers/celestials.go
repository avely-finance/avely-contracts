package helpers

import (
	"os"

	"github.com/Zilliqa/gozilliqa-sdk/v3/account"
	"github.com/joho/godotenv"

	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/ethereum/go-ethereum/accounts"
)

type Celestials struct {
	Owner       *account.Wallet
	Admin       *account.Wallet
	Verifier    *account.Wallet
	EvmDeployer *accounts.Account
}

var _sdk *core.AvelySDK

func NewCelestials(ownerKey, adminKey, verifierKey, evmDeployerKey string) *Celestials {
	owner := account.NewWallet()
	owner.AddByPrivateKey(ownerKey)

	admin := account.NewWallet()
	admin.AddByPrivateKey(adminKey)

	verifier := account.NewWallet()
	verifier.AddByPrivateKey(verifierKey)

	deployerAcc, _ := _sdk.Evm.AddAccountByPrivateKey(evmDeployerKey)

	return &Celestials{
		Owner:       owner,
		Admin:       admin,
		Verifier:    verifier,
		EvmDeployer: deployerAcc,
	}
}

func LoadCelestialsFromEnv(sdk *core.AvelySDK, chain string) *Celestials {
	_sdk = sdk
	path := ".env." + chain
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("WARNING! There is no '%s' file. Please, make sure you set up the correct ENV manually", path)
	}
	return NewCelestials(os.Getenv("OWNERKEY"), os.Getenv("ADMINKEY"), os.Getenv("VERIFIERKEY"),
		os.Getenv("EVMDEPLOYERKEY"))
}
