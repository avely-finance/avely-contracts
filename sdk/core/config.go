package core

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ApiConfig struct {
	HttpUrl string
}

type Config struct {
	Chain                    string
	GasLimit                 string
	Api                      ApiConfig
	ChainId                  int
	TxConfrimMaxAttempts     int
	TxConfirmIntervalSec     int
	TxRetryCount             int
	StZilSsnAddress          string
	StZilSsnRewardShare      string
	HolderInitialDelegateZil int
	SsnInitialDelegateZil    int
	ProtocolRewardsFee       int

	MultisigAddr string
	ASwapAddr    string
	ZproxyAddr   string
	GzilAddr     string
	ZimplAddr    string
	StZilAddr    string
	BufferAddrs  []string
	SsnAddrs     []string
	HolderAddr   string
	TreasuryAddr string
}

func NewConfig(configPath string, chain string) *Config {
	config := &Config{
		Chain: chain,
	}

	path := ".env." + chain
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("WARNING! There is no '%s' file. Please, make sure you set up the correct ENV manually", path)
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %w \n", err)
	}

	section := viper.Sub(chain)
	section.AutomaticEnv()
	section.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if section == nil { // Sub returns nil if the key cannot be found
		log.Fatalf("Chain %s not found in config", chain)
	}

	err = section.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Fatal error config file: %w \n", err)
	}
	return config
}
