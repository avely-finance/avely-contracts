package main

import (
	"flag"
	"strconv"

	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
)

var log *Log
var sdk *AvelySDK

func main() {
	chainPtr := flag.String("chain", "local", "chain")
	tagPtr := flag.String("tag", "default", "specific command")
	ssnPtr := flag.String("ssn", "", "ssn address")

	flag.Parse()

	tag := *tagPtr

	log = NewLog()
	config := NewConfig(*chainPtr)
	sdk = NewAvelySDK(*config)

	shortcuts := map[string]string{
		"stzilssn": config.StZilSsnAddress,
		"addr1":    config.Addr1,
		"addr2":    config.Addr2,
		"addr3":    config.Addr3,
		"admin":    config.Admin,
		"verifier": config.Verifier,
	}
	log.AddShortcuts(shortcuts)

	m, _ := RestoreMultisigContract(sdk, config.Owner, []string{}, 0)

	m.UpdateWallet(config.OwnerKey)

	stZilAddr := config.StZilAddr

	switch tag {
	case "SubmitSetHolderAddressTransaction":
		setHolder(m, stZilAddr, config.HolderAddr)
	case "SubmitChangeBuffersTransaction":
		changeBuffers(m, stZilAddr, config.BufferAddrs)
	case "SubmitAddSSNTransaction":
		ssnaddr := *ssnPtr
		if ssnaddr == "" {
			log.Fatal("SSN address empty")
		}
		addSSN(m, stZilAddr, ssnaddr)
	case "SubmitRemoveSSNTransaction":
		ssnaddr := *ssnPtr
		if ssnaddr == "" {
			log.Fatal("SSN address empty")
		}
		removeSSN(m, stZilAddr, ssnaddr)
	case "SubmitClaimOwnerTransaction":
		claimOwner(m, stZilAddr)
	case "SubmitChangeRewardsFeeTransaction":
		changeRewardsFee(m, stZilAddr, strconv.Itoa(sdk.Cfg.ProtocolRewardsFee))
	case "SubmitChangeTreasuryAddressTransaction":
		changeTreasuryAddress(m, stZilAddr, sdk.Cfg.TreasuryAddr)
	case "SubmitUnPauseInTransaction":
		unPauseIn(m, stZilAddr)
	case "SubmitUnPauseOutTransaction":
		unPauseOut(m, stZilAddr)
	case "SubmitUnPauseZrc2Transaction":
		unauseZrc2(m, stZilAddr)
	default:
		log.Fatal("Unknown tx tag")
	}

	log.Info("Done")
}

func setHolder(m *MultisigWallet, callee string, new_holder string) {
	check(m.SubmitSetHolderAddressTransaction(callee, new_holder))
}

func changeBuffers(m *MultisigWallet, callee string, buffers []string) {
	check(m.SubmitChangeBuffersTransaction(callee, buffers))
}

func addSSN(m *MultisigWallet, callee string, ssnaddr string) {
	check(m.SubmitAddSSNTransaction(callee, ssnaddr))
}

func removeSSN(m *MultisigWallet, callee string, ssnaddr string) {
	check(m.SubmitRemoveSSNTransaction(callee, ssnaddr))
}

func changeRewardsFee(m *MultisigWallet, callee string, value string) {
	check(m.SubmitChangeRewardsFeeTransaction(callee, value))
}

func changeTreasuryAddress(m *MultisigWallet, callee string, value string) {
	check(m.SubmitChangeTreasuryAddressTransaction(callee, value))
}

func claimOwner(m *MultisigWallet, callee string) {
	check(m.SubmitClaimOwnerTransaction(callee))
}

func unPauseIn(m *MultisigWallet, callee string) {
	check(m.SubmitUnPauseInTransaction(callee))
}

func unPauseOut(m *MultisigWallet, callee string) {
	check(m.SubmitUnPauseOutTransaction(callee))
}

func unauseZrc2(m *MultisigWallet, callee string) {
	check(m.SubmitUnPauseZrc2Transaction(callee))
}

func check(tx *transaction.Transaction, err error) {
	if err != nil {
		log.Error("Err: " + err.Error())
	} else {
		log.Info(tx)
	}
}
