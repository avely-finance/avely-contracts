package core

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ApiConfig struct {
	HttpUrl         string
	WebsocketUrl    string
	WebsocketSchema string
}

type SlackConfig struct {
	HookUrl  string
	LogLevel string
}

type Config struct {
	Chain                    string
	Api                      ApiConfig
	Slack                    SlackConfig
	ChainId                  int
	TxConfrimMaxAttempts     int
	TxConfirmIntervalSec     int
	StZilSsnAddress          string
	StZilSsnRewardShare      string
	HolderInitialDelegateZil int
	SsnInitialDelegateZil    int
	ProtocolRewardsFee       int
	Owner                    string
	OwnerKey                 string
	Admin                    string
	AdminKey                 string
	Addr1                    string
	Key1                     string
	Addr2                    string
	Key2                     string
	Addr3                    string
	Key3                     string
	Verifier                 string
	VerifierKey              string

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

func NewConfig(chain string) *Config {
	config := &Config{
		Chain: chain,
	}

	path := ".env." + chain
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("WARNING! There is no '%s' file. Please, make sure you set up the correct ENV manually", path)
	}

	viper.AddConfigPath(".")
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
