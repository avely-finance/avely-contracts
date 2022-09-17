package main

import (
	"flag"
	"reflect"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	. "github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/sirupsen/logrus"
)

type TreasuryAdminCli struct {
	sdk        *AvelySDK
	config     *Config
	treasury   *TreasuryContract
	multisig   *MultisigWallet
	isMultisig bool
	chain      string
	actor      string
	actorKey   string
	cmd        string
	recipient  string
	amount     uint
}

var log *Log

func main() {

	//init
	cli, err := NewTreasuryAdminCli()
	if err != nil {
		log.WithError(err).Fatal("Can't initialize Treasury CLI")
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

	//restore Treasury contract
	cli.treasury = cli.restoreTreasury()
	cli.treasury.UpdateWallet(cli.actorKey)

	//process command
	switch cli.cmd {
	case "print_state":
		cli.printState()
	//case "set_liquidity_fee":
	//	cli.setLiquidityFee()
	case "change_owner":
		cli.changeOwner()
	case "claim_owner":
		cli.claimOwner()
	case "withdraw":
		cli.Withdraw()
	default:
		log.WithFields(logrus.Fields{
			"chain":     cli.chain,
			"command":   cli.cmd,
			"recipient": cli.recipient,
		}).Fatal("Unknown command")
	}

	log.Info("Done")
}

func NewTreasuryAdminCli() (*TreasuryAdminCli, error) {
	chainPtr := flag.String("chain", "local", "chain")
	usrPtr := flag.String("usr", "owner", "an user ID/admin/owner")
	cmdPtr := flag.String("cmd", "default", "specific command")
	recipientPtr := flag.String("recipient", "", "recipient")
	amountPtr := flag.Uint("amount", 0, "amount, ZIL")

	flag.Parse()

	log = NewLog()
	config := NewConfig(*chainPtr)
	sdk := NewAvelySDK(*config)

	return &TreasuryAdminCli{
		sdk:        sdk,
		config:     config,
		chain:      *chainPtr,
		actor:      *usrPtr,
		actorKey:   getActorKey(sdk, *usrPtr),
		cmd:        *cmdPtr,
		recipient:  *recipientPtr,
		amount:     *amountPtr,
		isMultisig: false,
	}, nil
}

func (cli *TreasuryAdminCli) deploy() {
	//deploy is going from cli.config.Admin
	treasury, err := NewTreasuryContract(cli.sdk, cli.config.Owner)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain":       cli.chain,
			"command":     cli.cmd,
			"recipient":   cli.recipient,
			"owner":       cli.config.Owner,
			"is_multisig": cli.isMultisig,
		}).Fatal("Can't deploy Treasury contract")
	}
	log.WithFields(logrus.Fields{
		"chain":         cli.chain,
		"command":       cli.cmd,
		"recipient":     cli.recipient,
		"treasury_addr": treasury.Addr,
		"owner":         cli.config.Owner,
		"is_multisig":   cli.isMultisig,
	}).Info("Treasury contract deployed")
}

func (cli *TreasuryAdminCli) restoreTreasury() *TreasuryContract {
	treasury, err := RestoreTreasuryContract(cli.sdk, cli.config.TreasuryAddr, "")
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain":       cli.chain,
			"command":     cli.cmd,
			"recipient":   cli.recipient,
			"owner":       cli.config.Owner,
			"is_multisig": cli.isMultisig,
		}).Fatal("Can't restore Treasury contract")
	}
	log.WithFields(logrus.Fields{
		"chain":         cli.chain,
		"command":       cli.cmd,
		"recipient":     cli.recipient,
		"treasury_addr": treasury.Addr,
		"owner":         cli.config.Owner,
		"is_multisig":   cli.isMultisig,
	}).Info("Treasury contract restored")
	treasury.UpdateWallet(cli.config.OwnerKey)
	return treasury
}

func (cli *TreasuryAdminCli) restoreMultisig() *MultisigWallet {
	multisig, err := RestoreMultisigContract(cli.sdk, cli.config.Owner, []string{}, 0)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain":       cli.chain,
			"command":     cli.cmd,
			"recipient":   cli.recipient,
			"owner":       cli.config.Owner,
			"is_multisig": cli.isMultisig,
		}).Fatal("Can't restore Multisig contract")
	}
	log.WithFields(logrus.Fields{
		"chain":       cli.chain,
		"command":     cli.cmd,
		"recipient":   cli.recipient,
		"owner":       cli.config.Owner,
		"is_multisig": cli.isMultisig,
	}).Info("Multisig contract restored")
	//OwnerKey is key of user who will submit multisig transactions
	multisig.UpdateWallet(cli.config.OwnerKey)
	return multisig
}

func getActorKey(sdk *AvelySDK, usr string) string {
	var keyValue reflect.Value
	var key = ""
	if usr == "admin" {
		pr := reflect.ValueOf(sdk.Cfg)
		keyValue = reflect.Indirect(pr).FieldByName("AdminKey")
	} else if usr == "owner" {
		pr := reflect.ValueOf(sdk.Cfg)
		keyValue = reflect.Indirect(pr).FieldByName("OwnerKey")
	} else {
		pr := reflect.ValueOf(sdk.Cfg)
		keyValue = reflect.Indirect(pr).FieldByName("Key" + usr)
	}

	if keyValue.Kind() == reflect.Invalid {
		log.WithFields(logrus.Fields{
			"key": usr,
		}).Fatal("Can't get actor key")
	}

	key = keyValue.Interface().(string)

	return key
}

func (cli *TreasuryAdminCli) printState() {
	log.Info(cli.treasury.State())
}

/*func (cli *TreasuryAdminCli) setLiquidityFee() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitSetLiquidityFeeTransaction(cli.treasury.Addr, cli.recipient)
	} else {
		//address setup
		tx, err = cli.treasury.SetLiquidityFee(cli.recipient)
	}

	if err != nil {
		cli.logFatal("SetLiquidityFee error", err)
	}
	cli.logInfo("SetLiquidityFee succeed", tx)
}*/

func (cli *TreasuryAdminCli) changeOwner() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitChangeOwnerTransaction(cli.treasury.Addr, cli.recipient)
	} else {
		//address setup
		tx, err = cli.treasury.ChangeOwner(cli.recipient)
	}

	if err != nil {
		cli.logFatal("ChangeOwner error", err)
	}
	cli.logInfo("ChangeOwner succeed", tx)
}

func (cli *TreasuryAdminCli) Withdraw() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitWithdrawTransaction(cli.treasury.Addr, cli.recipient, ToQA(int(cli.amount)))
	} else {
		//address setup
		tx, err = cli.treasury.Withdraw(cli.recipient, ToQA(int(cli.amount)))
	}

	if err != nil {
		cli.logFatal("Withdraw error", err)
	}
	cli.logInfo("Withdraw succeed", tx)
}

func (cli *TreasuryAdminCli) claimOwner() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitClaimOwnerTransaction(cli.treasury.Addr)
	} else {
		//address setup
		tx, err = cli.treasury.ClaimOwner()
	}

	if err != nil {
		cli.logFatal("ClaimOwner error", err)
	}
	cli.logInfo("ClaimOwner succeed", tx)
}

func (cli *TreasuryAdminCli) addressIsContract(address string) bool {
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

func (cli *TreasuryAdminCli) logFatal(message string, err error) {
	b32, _ := bech32.ToBech32Address(cli.treasury.Addr)
	log.WithError(err).WithFields(logrus.Fields{
		"chain":             cli.chain,
		"actor":             cli.actor,
		"command":           cli.cmd,
		"recipient":         cli.recipient,
		"treasury_addr":     cli.treasury.Addr,
		"treasury_addr_b32": b32,
		"owner":             cli.config.Owner,
		"is_multisig":       cli.isMultisig,
	}).Fatal(message)
}

func (cli *TreasuryAdminCli) logInfo(message string, tx *transaction.Transaction) {
	b32, _ := bech32.ToBech32Address(cli.treasury.Addr)
	log.WithFields(logrus.Fields{
		"chain":             cli.chain,
		"actor":             cli.actor,
		"command":           cli.cmd,
		"recipient":         cli.recipient,
		"treasury_addr":     cli.treasury.Addr,
		"treasury_addr_b32": b32,
		"owner":             cli.config.Owner,
		"is_multisig":       cli.isMultisig,
		"txid":              tx.ID,
	}).Info(message)
}
