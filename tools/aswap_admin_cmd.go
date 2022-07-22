package main

import (
	"flag"

	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/sirupsen/logrus"
)

type ASwapCli struct {
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
	cli, err := NewASwapCli()
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
	default:
		log.WithFields(logrus.Fields{
			"chain":   cli.chain,
			"command": cli.cmd,
			"value":   cli.value1,
		}).Fatal("Unknown command")
	}

	log.Info("Done")
}

func NewASwapCli() (*ASwapCli, error) {
	chainPtr := flag.String("chain", "local", "chain")
	cmdPtr := flag.String("cmd", "default", "specific command")
	valuePtr := flag.String("value", "", "new value")

	flag.Parse()

	log = NewLog()
	config := NewConfig(*chainPtr)

	return &ASwapCli{
		sdk:        NewAvelySDK(*config),
		config:     config,
		chain:      *chainPtr,
		cmd:        *cmdPtr,
		value1:     *valuePtr,
		isMultisig: false,
	}, nil
}

func (cli *ASwapCli) deploy() {
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

func (cli *ASwapCli) restoreASwap() *ASwap {
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

func (cli *ASwapCli) restoreMultisig() *MultisigWallet {
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

func (cli *ASwapCli) printState() {
	log.Info(cli.aswap.State())
}

func (cli *ASwapCli) togglePause() {
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

func (cli *ASwapCli) setLiquidityFee() {
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

func (cli *ASwapCli) setTreasuryFee() {
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

func (cli *ASwapCli) setTreasuryAddress() {
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

func (cli *ASwapCli) changeOwner() {
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

func (cli *ASwapCli) addressIsContract(address string) bool {
	provider := provider2.NewProvider(cli.config.Api.HttpUrl)
	_, err := provider.GetSmartContractState(address[2:])
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain":         cli.chain,
			"command":       cli.cmd,
			"aswap_addr":    cli.aswap.Addr,
			"owner":         cli.config.Owner,
			"param_address": address,
		}).Debug("Address is not contract")
		return false
	}
	return true
}

func (cli *ASwapCli) logFatal(message string, err error) {
	log.WithError(err).WithFields(logrus.Fields{
		"chain":       cli.chain,
		"command":     cli.cmd,
		"value":       cli.value1,
		"aswap_addr":  cli.aswap.Addr,
		"owner":       cli.config.Owner,
		"is_multisig": cli.isMultisig,
	}).Fatal(message)
}

func (cli *ASwapCli) logInfo(message string, tx *transaction.Transaction) {
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
