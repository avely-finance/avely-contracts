package main

import (
	"flag"
	"reflect"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"

	//. "github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/sirupsen/logrus"
)

type SsnAdminCli struct {
	sdk        *AvelySDK
	config     *Config
	ssn        *SsnContract
	multisig   *MultisigWallet
	isMultisig bool
	chain      string
	actor      string
	actorKey   string
	cmd        string
	param1     string
}

var log *Log

func main() {

	//init
	cli, err := NewSsnAdminCli()
	if err != nil {
		log.WithError(err).Fatal("Can't initialize Ssn CLI")
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

	//restore Ssn contract
	cli.ssn = cli.restoreSsn()
	cli.ssn.UpdateWallet(cli.actorKey)

	//process command
	switch cli.cmd {
	case "print_state":
		cli.printState()
	case "change_owner":
		cli.changeOwner()
	case "claim_owner":
		cli.claimOwner()
	case "withdraw_comm":
		cli.WithdrawComm()
	case "update_comm":
		cli.UpdateComm()
	case "change_zproxy":
		cli.ChangeZproxy()
	case "update_receiving_addr":
		cli.UpdateReceivingAddr()
	default:
		log.WithFields(logrus.Fields{
			"chain":   cli.chain,
			"command": cli.cmd,
			"param1":  cli.param1,
		}).Fatal("Unknown command")
	}

	log.Info("Done")
}

func NewSsnAdminCli() (*SsnAdminCli, error) {
	chainPtr := flag.String("chain", "local", "chain")
	usrPtr := flag.String("usr", "owner", "an user ID/admin/owner")
	cmdPtr := flag.String("cmd", "default", "specific command")
	param1Ptr := flag.String("param1", "", "param1")

	flag.Parse()

	log = NewLog()
	config := NewConfig(*chainPtr)
	sdk := NewAvelySDK(*config)

	return &SsnAdminCli{
		sdk:        sdk,
		config:     config,
		chain:      *chainPtr,
		actor:      *usrPtr,
		actorKey:   getActorKey(sdk, *usrPtr),
		cmd:        *cmdPtr,
		param1:     *param1Ptr,
		isMultisig: false,
	}, nil
}

func (cli *SsnAdminCli) deploy() {
	//deploy is going from cli.config.Admin
	ssn, err := NewSsnContract(cli.sdk, cli.config.Owner, cli.config.ZproxyAddr)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain":       cli.chain,
			"command":     cli.cmd,
			"param1":      cli.param1,
			"owner":       cli.config.Owner,
			"is_multisig": cli.isMultisig,
		}).Fatal("Can't deploy Ssn contract")
	}
	log.WithFields(logrus.Fields{
		"chain":       cli.chain,
		"command":     cli.cmd,
		"param1":      cli.param1,
		"ssn_addr":    ssn.Addr,
		"owner":       cli.config.Owner,
		"is_multisig": cli.isMultisig,
	}).Info("Ssn contract deployed")
}

func (cli *SsnAdminCli) restoreSsn() *SsnContract {
	ssn, err := RestoreSsnContract(cli.sdk, cli.config.StZilSsnAddress, "", "")
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain":       cli.chain,
			"command":     cli.cmd,
			"param1":      cli.param1,
			"owner":       cli.config.Owner,
			"is_multisig": cli.isMultisig,
		}).Fatal("Can't restore Ssn contract")
	}
	log.WithFields(logrus.Fields{
		"chain":       cli.chain,
		"command":     cli.cmd,
		"param1":      cli.param1,
		"ssn_addr":    ssn.Addr,
		"owner":       cli.config.Owner,
		"is_multisig": cli.isMultisig,
	}).Info("Ssn contract restored")
	ssn.UpdateWallet(cli.config.OwnerKey)
	return ssn
}

func (cli *SsnAdminCli) restoreMultisig() *MultisigWallet {
	multisig, err := RestoreMultisigContract(cli.sdk, cli.config.Owner, []string{}, 0)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"chain":       cli.chain,
			"command":     cli.cmd,
			"param1":      cli.param1,
			"owner":       cli.config.Owner,
			"is_multisig": cli.isMultisig,
		}).Fatal("Can't restore Multisig contract")
	}
	log.WithFields(logrus.Fields{
		"chain":       cli.chain,
		"command":     cli.cmd,
		"param1":      cli.param1,
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

func (cli *SsnAdminCli) printState() {
	log.Info(cli.ssn.State())
}

func (cli *SsnAdminCli) changeOwner() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitChangeOwnerTransaction(cli.ssn.Addr, cli.param1)
	} else {
		//address setup
		tx, err = cli.ssn.ChangeOwner(cli.param1)
	}

	if err != nil {
		cli.logFatal("ChangeOwner error", err)
	}
	cli.logInfo("ChangeOwner succeed", tx)
}

func (cli *SsnAdminCli) WithdrawComm() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitWithdrawCommTransaction(cli.ssn.Addr)
	} else {
		//address setup
		tx, err = cli.ssn.WithdrawComm()
	}

	if err != nil {
		cli.logFatal("WithdrawComm error", err)
	}
	cli.logInfo("WithdrawComm succeed", tx)
}

func (cli *SsnAdminCli) UpdateComm() {
	var tx *transaction.Transaction
	var err error
	new_comm, _ := strconv.Atoi(cli.param1)
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitUpdateCommTransaction(cli.ssn.Addr, new_comm)
	} else {
		//address setup
		tx, err = cli.ssn.UpdateComm(new_comm)
	}

	if err != nil {
		cli.logFatal("UpdateComm error", err)
	}
	cli.logInfo("UpdateComm succeed", tx)
}

func (cli *SsnAdminCli) ChangeZproxy() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitChangeZproxyTransaction(cli.ssn.Addr, cli.param1)
	} else {
		//address setup
		tx, err = cli.ssn.ChangeZproxy(cli.param1)
	}

	if err != nil {
		cli.logFatal("ChangeZproxy error", err)
	}
	cli.logInfo("ChangeZproxy succeed", tx)
}

func (cli *SsnAdminCli) UpdateReceivingAddr() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitUpdateReceivingAddrTransaction(cli.ssn.Addr, cli.param1)
	} else {
		//address setup
		tx, err = cli.ssn.UpdateReceivingAddr(cli.param1)
	}

	if err != nil {
		cli.logFatal("UpdateReceivingAddr error", err)
	}
	cli.logInfo("UpdateReceivingAddr succeed", tx)
}

func (cli *SsnAdminCli) claimOwner() {
	var tx *transaction.Transaction
	var err error
	if cli.isMultisig {
		//multisig setup
		tx, err = cli.multisig.SubmitClaimOwnerTransaction(cli.ssn.Addr)
	} else {
		//address setup
		tx, err = cli.ssn.ClaimOwner()
	}

	if err != nil {
		cli.logFatal("ClaimOwner error", err)
	}
	cli.logInfo("ClaimOwner succeed", tx)
}

func (cli *SsnAdminCli) addressIsContract(address string) bool {
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

func (cli *SsnAdminCli) logFatal(message string, err error) {
	b32, _ := bech32.ToBech32Address(cli.ssn.Addr)
	log.WithError(err).WithFields(logrus.Fields{
		"chain":        cli.chain,
		"actor":        cli.actor,
		"command":      cli.cmd,
		"param1":       cli.param1,
		"ssn_addr":     cli.ssn.Addr,
		"ssn_addr_b32": b32,
		"owner":        cli.config.Owner,
		"is_multisig":  cli.isMultisig,
	}).Fatal(message)
}

func (cli *SsnAdminCli) logInfo(message string, tx *transaction.Transaction) {
	b32, _ := bech32.ToBech32Address(cli.ssn.Addr)
	log.WithFields(logrus.Fields{
		"chain":        cli.chain,
		"actor":        cli.actor,
		"command":      cli.cmd,
		"param1":       cli.param1,
		"ssn_addr":     cli.ssn.Addr,
		"ssn_addr_b32": b32,
		"owner":        cli.config.Owner,
		"is_multisig":  cli.isMultisig,
		"txid":         tx.ID,
	}).Info(message)
}
