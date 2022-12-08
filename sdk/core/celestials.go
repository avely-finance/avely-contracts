package core

import (
	"log"
	"os"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/joho/godotenv"
)

type Celestials struct {
	Owner    *account.Wallet
	Admin    *account.Wallet
	Verifier *account.Wallet
}

func NewCelestials(ownerKey, adminKey, verifierKey string) *Celestials {
	owner := account.NewWallet()
	owner.AddByPrivateKey(ownerKey)

	admin := account.NewWallet()
	admin.AddByPrivateKey(adminKey)

	verifier := account.NewWallet()
	verifier.AddByPrivateKey(verifierKey)

	return &Celestials{
		Owner:    owner,
		Admin:    admin,
		Verifier: verifier,
	}
}

func LoadCelestialsFromEnv(chain string) *Celestials {
	path := ".env." + chain
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("WARNING! There is no '%s' file. Please, make sure you set up the correct ENV manually", path)
	}

	return NewCelestials(os.Getenv("OWNERKEY"), os.Getenv("ADMINKEY"), os.Getenv("VERIFIERKEY"))
}
