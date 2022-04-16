package actions

import (
	"errors"
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
	return &AdminActions{log: log}
}

func (a *AdminActions) SetTestMode(mode bool) {
	a.testMode = mode
	if a.testMode {
		a.TxLogMap = make(map[string]TxLog)
	} else {
		a.TxLogMap = nil
	}
}

func (a *AdminActions) SaveTx(step string, tx *transaction.Transaction, err error) (*transaction.Transaction, error) {
	if a.testMode {
		a.TxLogMap[step] = TxLog{tx, err}
	}
	return tx, err
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

	buffers := p.Azil.GetDrainedBuffers()
	if lastDrained, ok := buffers[bufferToDrain]; !ok {
		a.log.Info("Buffer is never drained; Let's do this first time")
	} else if lastDrained.Int() >= int64(lrc) {
		a.log.Debug("No need to drain buffer")
		return nil
	}

	ssnlist := p.Azil.GetSsnWhitelist()

	//claim rewards from holder
	for _, ssn := range ssnlist {
		tx, err := p.Azil.ClaimRewardsHolder(ssn)
		a.SaveTx("ClaimRewardsHolder_"+ssn, tx, err)
		if err != nil {
			a.log.WithFields(logrus.Fields{"tx": tx.ID, "ssn_address": ssn, "error": tx.Receipt}).Error("ClaimRewardsHolder failed")
			return errors.New("Buffer drain is failed at ClaimRewardsHolder step")
		} else {
			a.log.WithFields(logrus.Fields{"tx": tx.ID, "ssn_address": ssn}).Info("ClaimRewardsHolder success")
		}
	}

	//claim rewards from buffer
	for _, ssn := range ssnlist {
		tx, err := p.Azil.ClaimRewardsBuffer(bufferToDrain, ssn)
		a.SaveTx("ClaimRewardsBuffer_"+ssn, tx, err)
		if err != nil {
			a.log.WithFields(logrus.Fields{
				"tx":             tx.ID,
				"buffer_address": bufferToDrain,
				"ssn_address":    ssn,
				"error":          tx.Receipt,
			}).Error("ClaimRewardsBuffer failed")
			return errors.New("Buffer drain is failed at ClaimRewardsBuffer step")
		} else {
			a.log.WithFields(logrus.Fields{
				"tx":             tx.ID,
				"buffer_address": bufferToDrain,
				"ssn_address":    ssn,
			}).Info("ClaimRewardsBuffer success")
		}
	}

	//transfer stake from buffer to holder
	tx, err := p.Azil.ConsolidateInHolder(bufferToDrain)
	a.SaveTx("ConsolidateInHolder", tx, err)
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

	for fromSsn, amount := range mapSsnAmount {
		amountStr := amount.String()
		ssnForInput := p.GetSsnAddressForInput()
		if showOnly {
			a.log.WithFields(logrus.Fields{
				"from_ssn": fromSsn,
				"to_ssn":   ssnForInput,
				"amount":   amountStr,
			}).Debug("Need to call Azil.ChownStakeReDelegate transition")
		} else if tx, err := p.Azil.ChownStakeReDelegate(fromSsn, amountStr); err != nil {
			a.log.WithFields(logrus.Fields{
				"from_ssn": fromSsn,
				"to_ssn":   ssnForInput,
				"amount":   amountStr,
				"tx":       tx.ID,
				"error":    tx.Receipt,
			}).Error("ChownStakeReDelegate failed")
			a.SaveTx("ChownStakeReDelegate_"+fromSsn, tx, err)
			return errors.New("ChownStakeReDelegate failed")
		} else {
			a.log.WithFields(logrus.Fields{
				"from_ssn": fromSsn,
				"to_ssn":   ssnForInput,
				"amount":   amountStr,
				"tx":       tx.ID,
			}).Info("ChownStakeReDelegate OK")
			a.SaveTx("ChownStakeReDelegate_"+fromSsn, tx, err)
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
		tx, err := p.Azil.ChownStakeConfirmSwap(initiator)
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
	autorestakeamount := p.Azil.GetAutorestakeAmount()

	if autorestakeamount.Cmp(big.NewInt(0)) == 0 { // autorestakeamount == 0
		a.log.Info("Nothing to autorestake")
		return nil
	}

	if autorestakeamount.Cmp(big.NewInt(10000000000000)) < 1 { // autorestakeamount <= 10 ZIL
		a.log.Info("Autorestake is lower than min delegate amount. " + autorestakeamount.String())
		return nil
	}

	priceBefore := p.Azil.GetAzilPrice().String()
	tx, err := p.Azil.PerformAutoRestake()

	if err != nil {
		a.log.WithFields(logrus.Fields{"txid": tx.ID, "error": tx.Receipt}).Error("AutoRestake failed")

		return errors.New("AutoRestake failed")
	}

	priceAfter := p.Azil.GetAzilPrice().String()

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
	tx, err := p.Azil.ClaimWithdrawal(blocksStr)

	if err != nil {
		a.log.WithFields(logrus.Fields{"tx": tx.ID, "error": tx.Receipt}).Error("Withdrawals claim failed")

		return errors.New("Withdrawals claim failed")
	} else {
		a.log.WithFields(logrus.Fields{"tx": tx.ID}).Info("Withdrawals successfully claimed")
	}

	return nil
}
