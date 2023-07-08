package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/sirupsen/logrus"
)

type TxLog struct {
	Tx  *transaction.Transaction
	Err error
}

type AdminActions struct {
	log      *core.Log
	testMode bool
	TxLogMap map[string]TxLog
}

func NewAdminActions(log *core.Log) *AdminActions {
	return &AdminActions{
		log:      log,
		testMode: false,
		TxLogMap: make(map[string]TxLog),
	}
}

func (a *AdminActions) TxLogMode(mode bool) {
	a.testMode = mode
}

func (a *AdminActions) TxLog(step string, tx *transaction.Transaction, err error) (*transaction.Transaction, error) {
	if a.testMode {
		a.TxLogMap[step] = TxLog{tx, err}
	}
	return tx, err
}

func (a *AdminActions) TxLogClear() {
	a.TxLogMap = make(map[string]TxLog)
}

func (a *AdminActions) HasTxError(txn *transaction.Transaction, errorCode string) bool {
	receipt, _ := json.Marshal(txn.Receipt)
	receiptStr := string(receipt)
	errorMessage := fmt.Sprintf("Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 %s))])", errorCode)
	return strings.Contains(receiptStr, errorMessage)
}

func (a *AdminActions) DrainBufferAuto(p *Protocol) error {
	lrc := p.Zimpl.GetLastRewardCycle()
	bufferToDrain := p.GetBufferToDrain()
	return a.DrainBuffer(p, lrc, bufferToDrain.Addr)
}

func (a *AdminActions) DrainBufferByCycle(p *Protocol, lrc int) error {
	bufferToDrain := p.GetBufferToDrain()
	return a.DrainBuffer(p, lrc, bufferToDrain.Addr)
}

func (a *AdminActions) DrainBuffer(p *Protocol, lrc int, bufferToDrain string) error {
	buffers := p.StZIL.GetDrainedBuffers()
	if lastDrained, ok := buffers[bufferToDrain]; !ok {
		a.log.Info("Buffer is never drained; Let's do this first time")
	} else if lastDrained.Int() >= int64(lrc) {
		a.log.Debug("No need to drain buffer")
		return nil
	}

	ssnlist := p.StZIL.Sdk.Cfg.SsnAddrs

	//claim rewards from holder
	for _, ssn := range ssnlist {
		txCall := func() (*transaction.Transaction, error) { return p.StZIL.ClaimRewards(p.Holder.Addr, ssn) }
		tx, err := retryTx(p.StZIL.Sdk.Cfg.TxRetryCount, txCall)

		fields := logrus.Fields{
			"tx":             tx.ID,
			"holder_address": p.Holder.Addr,
			"ssn_address":    ssn,
		}
		a.TxLog("ClaimRewardsHolder_"+ssn, tx, err)
		if err == nil {
			a.log.WithFields(fields).Info("ClaimRewards Holder success")
		} else {
			fields["error"] = tx.Receipt
			a.log.WithFields(fields).Error("ClaimRewards Holder error")
		}
	}

	//claim rewards from buffer
	deposits := p.Zimpl.GetDepositAmtDeleg(bufferToDrain)
	bufferDeposits := p.Zimpl.GetBufferAmtDeleg(bufferToDrain)

	for _, ssn := range ssnlist {
		_, notNull := deposits[strings.ToLower(ssn)]
		_, bufferNotNull := bufferDeposits[strings.ToLower(ssn)]

		if notNull || bufferNotNull {
			txCall := func() (*transaction.Transaction, error) { return p.StZIL.ClaimRewards(bufferToDrain, ssn) }
			tx, err := retryTx(p.StZIL.Sdk.Cfg.TxRetryCount, txCall)

			fields := logrus.Fields{
				"tx":             tx.ID,
				"buffer_address": bufferToDrain,
				"ssn_address":    ssn,
			}
			a.TxLog("ClaimRewardsBuffer_"+ssn, tx, err)
			if err == nil {
				a.log.WithFields(fields).Info("ClaimRewards Buffer success")
			} else {
				fields["error"] = tx.Receipt
				a.log.WithFields(fields).Error("ClaimRewards Buffer error")
			}
		} else {
			a.log.WithFields(logrus.Fields{"buffer_address": bufferToDrain, "ssn": ssn}).Info("Was skipped because does not have stake")
		}
	}

	//transfer stake from buffer to holder
	txCall := func() (*transaction.Transaction, error) { return p.StZIL.ConsolidateInHolder(bufferToDrain) }
	tx, err := retryTx(p.StZIL.Sdk.Cfg.TxRetryCount, txCall)

	a.TxLog("ConsolidateInHolder", tx, err)
	if err != nil {
		a.log.WithFields(logrus.Fields{
			"tx":             tx.ID,
			"buffer_address": bufferToDrain,
			"error":          tx.Receipt,
		}).Error("ConsolidateInHolder failed")
		return errors.New("Buffer drain is failed at ConsolidateInHolder step")
	}

	a.log.WithFields(logrus.Fields{
		"buffer_address": bufferToDrain,
		"tx":             tx.ID,
	}).Info("ConsolidateInHolder Success; buffer successfully drained")

	return nil
}

func (a *AdminActions) ChownStakeReDelegate(p *Protocol, showOnly bool) error {
	activeBuffer := p.GetActiveBuffer()
	a.log.WithFields(logrus.Fields{"active_buffer": activeBuffer.Addr}).Info("Active Buffer")

	mapSsnAmount := p.Zimpl.GetDepositAmtDeleg(activeBuffer.Addr)
	if 0 == len(mapSsnAmount) {
		a.log.Info("Buffer is empty, nothing to redelegate.")
		return nil
	}

	ssnlist := p.StZIL.Sdk.Cfg.SsnAddrs

	for fromSsn, amount := range mapSsnAmount {
		if contains(ssnlist, fromSsn) {
			a.log.WithFields(logrus.Fields{"from_ssn": fromSsn, "amount": amount}).Info("SSN is in the list of whitelisted SSNs")
			continue
		}

		amountStr := amount.String()
		ssnForInput := p.GetSsnAddressForInput()
		fields := logrus.Fields{
			"from_ssn": fromSsn,
			"to_ssn":   ssnForInput,
			"amount":   amountStr,
		}

		txCall := func() (*transaction.Transaction, error) { return p.StZIL.ChownStakeReDelegate(fromSsn, amountStr) }

		if showOnly {
			a.log.WithFields(fields).Debug("Need to call StZIL.ChownStakeReDelegate transition")
		} else if tx, err := retryTx(p.StZIL.Sdk.Cfg.TxRetryCount, txCall); err == nil {
			a.TxLog("ChownStakeReDelegate_"+fromSsn, tx, err)
			fields["tx"] = tx.ID
			a.log.WithFields(fields).Info("ChownStakeReDelegate OK")
		} else {
			a.TxLog("ChownStakeReDelegate_"+fromSsn, tx, err)
			fields["tx"] = tx.ID
			fields["error"] = tx.Receipt
			if a.HasTxError(tx, p.StZIL.ErrorCode("DelegDoesNotExistAtSSN")) {
				a.log.WithFields(fields).Warning("ChownStakeReDelegate reported DelegDoesNotExistAtSSN error")
			} else {
				a.log.WithFields(fields).Error("ChownStakeReDelegate failed")
				return errors.New("ChownStakeReDelegate failed")
			}
		}
	}

	return nil
}

func (a *AdminActions) ProcessSwapRequests(p *Protocol, bufferOffset int) error {
	//bufferOffset=0 -> currentBuffer; bufferOffset=1 -> nextBuffer
	buffer := p.GetBufferByOffset(bufferOffset).Addr
	swapRequests := p.GetSwapRequestsForBuffer(buffer)
	a.log.WithFields(logrus.Fields{
		"swap_requests_count": len(swapRequests),
		"buffer":              buffer,
		"bufferOffset":        bufferOffset,
	}).Debug("Swap request(s) found")
	errCnt := 0
	okCnt := 0
	for _, initiator := range swapRequests {
		tx, err := p.StZIL.ChownStakeConfirmSwap(initiator)
		if err != nil {
			a.log.WithFields(logrus.Fields{"initiator": initiator, "txid": tx.ID, "error": tx.Receipt}).Error("Can't process swap")
			errCnt++
		} else {
			a.log.WithFields(logrus.Fields{"initiator": initiator, "txid": tx.ID}).Info("Swap processed")
			okCnt++
		}
	}
	a.log.WithFields(logrus.Fields{"processed_swaps_count": okCnt, "errors_count": errCnt}).Debug("ProcessSwapRequests completed")

	if errCnt > 0 {
		return errors.New("ProcessSwapRequests has failed items")
	}
	return nil
}

func (a *AdminActions) AutoRestake(p *Protocol) error {
	autorestakeamount := p.StZIL.GetAutorestakeAmount()

	if autorestakeamount.Cmp(big.NewInt(0)) == 0 { // autorestakeamount == 0
		a.log.Info("Nothing to autorestake")
		return nil
	}

	minDelegStake := p.StZIL.GetMinDelegStake()
	minDelegStakeBI, _ := big.NewInt(0).SetString(minDelegStake, 10)
	if autorestakeamount.Cmp(minDelegStakeBI) < 0 { // autorestakeamount < 10 ZIL (default value)
		a.log.WithFields(logrus.Fields{
			"autorestake":   autorestakeamount.String(),
			"mindelegstake": minDelegStake,
		}).Info("Autorestake is lower than min delegate amount")
		return nil
	}

	priceBefore := p.StZIL.GetStZilPrice().String()

	txCall := func() (*transaction.Transaction, error) { return p.StZIL.PerformAutoRestake() }
	tx, err := retryTx(p.StZIL.Sdk.Cfg.TxRetryCount, txCall)

	if err != nil {
		a.log.WithFields(logrus.Fields{"txid": tx.ID, "error": tx.Receipt}).Error("AutoRestake failed")

		return errors.New("AutoRestake failed")
	}

	priceAfter := p.StZIL.GetStZilPrice().String()

	a.log.WithFields(logrus.Fields{
		"tx":           tx.ID,
		"amount":       autorestakeamount.String(),
		"price_before": priceBefore,
		"price_after":  priceAfter,
	}).Info("AutoRestake is successfully completed")

	return nil
}

func (a *AdminActions) ShowClaimWithdrawal(p *Protocol) {
	blocks := p.GetClaimWithdrawalBlocks()
	if len(blocks) > 0 {
		blocksStr := utils.ArrayItoa(blocks)
		a.log.Debug("Blocks with unbonded withdrawals: " + strings.Join(blocksStr, ", "))
	} else {
		a.log.Debug("Blocks with unbonded withdrawals not found.")
	}
}

func (a *AdminActions) ClaimWithdrawal(p *Protocol) error {
	blocks := p.GetClaimWithdrawalBlocks()
	cnt := len(blocks)
	if cnt == 0 {
		a.log.Debug("There are no blocks with unbonded withdrawals.")
		return nil
	}
	blocksStr := utils.ArrayItoa(blocks)
	a.log.WithFields(logrus.Fields{"blocks_count": cnt, "blocks_list": strings.Join(blocksStr, ", ")}).Debug("Found blocks with unbonded withdrawals")
	tx, err := p.StZIL.ClaimWithdrawal(blocksStr)

	if err != nil {
		a.log.WithFields(logrus.Fields{"tx": tx.ID, "error": tx.Receipt}).Error("Withdrawals claim failed")

		return errors.New("Withdrawals claim failed")
	} else {
		a.log.WithFields(logrus.Fields{"tx": tx.ID}).Info("Withdrawals successfully claimed")
	}

	return nil
}

func retryTx(count int, txCall func() (*transaction.Transaction, error)) (*transaction.Transaction, error) {
	var tx *transaction.Transaction
	var err error

	for i := 0; i < count; i++ {
		tx, err = txCall()

		if err == nil {
			break
		}
	}

	return tx, err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
