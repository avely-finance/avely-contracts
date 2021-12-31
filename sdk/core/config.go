package core

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Chain                    string
	ApiUrl                   string
	TxConfrimMaxAttempts     int
	TxConfirmIntervalSec     int
	AzilSsnAddress           string
	AzilSsnRewardShare       string
	HolderInitialDelegateZil int
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

	ZproxyAddr  string
	GzilAddr    string
	ZimplAddr   string
	AproxyAddr  string
	AzilAddr    string
	BufferAddrs []string
	HolderAddr  string
}

func NewConfig(chain string) *Config {
	config := &Config{
		Chain: chain,
	}

	path := ".env." + chain
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Error loading %s file", path)
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %w \n", err)
	}

	section := viper.Sub(chain)
	section.AutomaticEnv()
	if section == nil { // Sub returns nil if the key cannot be found
		log.Fatalf("Chain %s not found in config", chain)
	}

	err = section.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Fatal error config file: %w \n", err)
	}

	return config
}
