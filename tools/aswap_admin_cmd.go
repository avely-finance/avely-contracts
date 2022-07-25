package main

import (
	"flag"
	"strings"

	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/sirupsen/logrus"
)

type ASwapAdminCli struct {
	sdk        *AvelySDK
	config     *Config
	chain      string
	cmd        string
	value1     string
	aswap      *ASwap
	multisig   *MultisigWallet
	isMultisig bool
}

var log *Log

func main() {

	//init
	cli, err := NewASwapAdminCli()
	if err != nil {
		log.WithError(err).Fatal("Can't initialize ASwap CLI")
	}

	//is owner multisig?
	if cli.addressIsContract(cli.config.Owner) {
		cli.isMultisig = true
	}

	//process deploy command
	if cli.cmd == "deploy" {
		cli.deploy()
		log.Info("Done")
		return
	}

	//restore Multisig contact
	if cli.isMultisig {
		cli.multisig = cli.restoreMultisig()
	}

	//restore ASwap contract
	cli.aswap = cli.restoreASwap()

	//process command
	switch cli.cmd {
	case "print_state":
		cli.printState()
	case "toggle_pause":
		cli.togglePause()
	case "set_liquidity_fee":
		cli.setLiquidityFee()
	case "set_treasury_fee":
		cli.setTreasuryFee()
	case "set_treasury_address":
		cli.setTreasuryAddress()
	case "change_owner":
		cli.changeOwner()
	case "claim_owner":
		cli.claimOwner()
	default:
		log.WithFields(logrus.Fields{
			"chain":   cli.chain,
			"command": cli.cmd,
			"value":   cli.value1,
		}).Fatal("Unknown command")
	}

	log.Info("Done")
}

func NewASwapAdminCli() (*ASwapAdminCli, error) {
	chainPtr := flag.String("chain", "local", "chain")
	cmdPtr := flag.String("cmd", "default", "specific command")
	valuePtr := flag.String("value", "", "new value")

	flag.Parse()

	log = NewLog()
	config := NewConfig(*chainPtr)

	return &ASwapAdminCli{
		sdk:        NewAvelySDK(*config),
		config:     config,
		chain:      *chainPtr,
		cmd:        *cmdPtr,
		value1:     *valuePtr,
		isMultisig: false,
	}, nil
}

func (cli *ASwapAdminCli) deploy() {
	//deploy is going from cli.config.Admin
	aswap, err := NewASwap(cli.sdk, cli.config.Owner)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain":       cli.chain,
			"command":     cli.cmd,
			"value":       cli.value1,
			"owner":       cli.config.Owner,
			"is_multisig": cli.isMultisig,
		}).Fatal("Can't deploy ASwap contract")
	}
	log.WithFields(logrus.Fields{
		"chain":       cli.chain,
		"command":     cli.cmd,
		"value":       cli.value1,
		"aswap_addr":  aswap.Addr,
		"owner":       cli.config.Owner,
		"is_multisig": cli.isMultisig,
	}).Info("ASwap contract deployed")
}

func (cli *ASwapAdminCli) restoreASwap() *ASwap {
	aswap, err := RestoreASwap(cli.sdk, cli.config.ASwapAddr, "")
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain":       cli.chain,
			"command":     cli.cmd,
			"value":       cli.value1,
			"owner":       cli.config.Owner,
			"is_multisig": cli.isMultisig,
		}).Fatal("Can't restore ASwap contract")
	}
	log.WithFields(logrus.Fields{
		"chain":       cli.chain,
		"command":     cli.cmd,
		"value":       cli.value1,
		"aswap_addr":  aswap.Addr,
		"owner":       cli.config.Owner,
		"is_multisig": cli.isMultisig,
	}).Info("ASwap contract restored")
	aswap.UpdateWallet(cli.config.OwnerKey)
	return aswap
}

func (cli *ASwapAdminCli) restoreMultisig() *MultisigWallet {
	multisig, err := RestoreMultisigContract(cli.sdk, cli.config.Owner, []string{}, 0)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain":       cli.chain,
			"command":     cli.cmd,
			"value":       cli.value1,
			"owner":       cli.config.Owner,
			"is_multisig": cli.isMultisig,
		}).Fatal("Can't restore Multisig contract")
	}
	log.WithFields(logrus.Fields{
		"chain":       cli.chain,
		"command":     cli.cmd,
		"value":       cli.value1,
		"owner":       cli.config.Owner,
		"is_multisig": cli.isMultisig,
	}).Info("Multisig contract restored")
	//OwnerKey is key of user who will submit multisig transactions
	multisig.UpdateWallet(cli.config.OwnerKey)
	return multisig
}

func (cli *ASwapAdminCli) printState() {
	log.Info(cli.aswap.State())
}

func (cli *ASwapAdminCli) togglePause() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitTogglePauseTransaction(cli.aswap.Addr)
	} else {
		//address setup
		tx, err = cli.aswap.TogglePause()
	}

	if err != nil {
		cli.logFatal("TogglePause error", err)
	}
	cli.logInfo("TogglePause succeed", tx)
}

func (cli *ASwapAdminCli) setLiquidityFee() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitSetLiquidityFeeTransaction(cli.aswap.Addr, cli.value1)
	} else {
		//address setup
		tx, err = cli.aswap.SetLiquidityFee(cli.value1)
	}

	if err != nil {
		cli.logFatal("SetLiquidityFee error", err)
	}
	cli.logInfo("SetLiquidityFee succeed", tx)
}

func (cli *ASwapAdminCli) setTreasuryFee() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitSetTreasuryFeeTransaction(cli.aswap.Addr, cli.value1)
	} else {
		//address setup
		tx, err = cli.aswap.SetTreasuryFee(cli.value1)
	}

	if err != nil {
		cli.logFatal("SetTreasuryFee error", err)
	}
	cli.logInfo("SetTreasuryFee succeed", tx)
}

func (cli *ASwapAdminCli) setTreasuryAddress() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitSetTreasuryAddressTransaction(cli.aswap.Addr, cli.value1)
	} else {
		//address setup
		tx, err = cli.aswap.SetTreasuryAddress(cli.value1)
	}

	if err != nil {
		cli.logFatal("SetTreasuryAddress error", err)
	}
	cli.logInfo("SetTreasuryAddress succeed", tx)
}

func (cli *ASwapAdminCli) changeOwner() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitChangeOwnerTransaction(cli.aswap.Addr, cli.value1)
	} else {
		//address setup
		tx, err = cli.aswap.ChangeOwner(cli.value1)
	}

	if err != nil {
		cli.logFatal("ChangeOwner error", err)
	}
	cli.logInfo("ChangeOwner succeed", tx)
}

func (cli *ASwapAdminCli) claimOwner() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitClaimOwnerTransaction(cli.aswap.Addr)
	} else {
		//address setup
		tx, err = cli.aswap.ClaimOwner()
	}

	if err != nil {
		cli.logFatal("ClaimOwner error", err)
	}
	cli.logInfo("ClaimOwner succeed", tx)
}

func (cli *ASwapAdminCli) addressIsContract(address string) bool {
	provider := provider2.NewProvider(cli.config.Api.HttpUrl)
	result, err := provider.GetSmartContractState(address[2:])
	if err != nil {
		//may be network error
		log.WithError(err).WithFields(logrus.Fields{
			"chain":         cli.chain,
			"command":       cli.cmd,
			"owner":         cli.config.Owner,
			"param_address": address,
		}).Fatal("Can't get owner address/contract type")
	} else if result.Error == nil {
		//there is no error, state fetched
		log.WithFields(logrus.Fields{
			"chain":   cli.chain,
			"command": cli.cmd,
			"owner":   cli.config.Owner,
			"address": address,
		}).Debug("Address is contract")
		return true

	}
	//https://dev.zilliqa.com/docs/apis/api-contract-get-smartcontract-state
	//-5:Address not contract address
	msg := strings.ToLower(result.Error.Message)
	if -1 != strings.Index(msg, "not contract") {
		log.WithFields(logrus.Fields{
			"chain":     cli.chain,
			"command":   cli.cmd,
			"owner":     cli.config.Owner,
			"address":   address,
			"rpc_error": result.Error.Message,
		}).Debug("Address is not contract")
		return false
	}

	log.WithFields(logrus.Fields{
		"chain":     cli.chain,
		"command":   cli.cmd,
		"owner":     cli.config.Owner,
		"address":   address,
		"rpc_error": result.Error.Message,
	}).Fatal("Can't get owner address/contract type")
	return false
}

func (cli *ASwapAdminCli) logFatal(message string, err error) {
	log.WithError(err).WithFields(logrus.Fields{
		"chain":       cli.chain,
		"command":     cli.cmd,
		"value":       cli.value1,
		"aswap_addr":  cli.aswap.Addr,
		"owner":       cli.config.Owner,
		"is_multisig": cli.isMultisig,
	}).Fatal(message)
}

func (cli *ASwapAdminCli) logInfo(message string, tx *transaction.Transaction) {
	log.WithFields(logrus.Fields{
		"chain":       cli.chain,
		"command":     cli.cmd,
		"value":       cli.value1,
		"aswap_addr":  cli.aswap.Addr,
		"owner":       cli.config.Owner,
		"is_multisig": cli.isMultisig,
		"txid":        tx.ID,
	}).Info(message)
}
