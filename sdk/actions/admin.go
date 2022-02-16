package actions

import (
	"errors"
	"math/big"
	"strings"

	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/sirupsen/logrus"
)

type AdminActions struct {
	log *core.Log
}

func NewAdminActions(log *core.Log) *AdminActions {
	return &AdminActions{log: log}
}

func (a *AdminActions) DrainBuffer(p *Protocol, lrc int) error {
	bufferToDrain := p.GetBufferToDrain()

	buffers := p.Azil.GetDrainedBuffers()
	needDrain := false

	if lastDrained, ok := buffers[bufferToDrain.Addr]; ok {
		if lastDrained.Int() != int64(lrc) {
			needDrain = true
		}
	} else {
		a.log.Info("Buffer is never drained; Let's do this first time")

		needDrain = true
	}

	if needDrain {
		tx, err := p.Azil.DrainBuffer(bufferToDrain.Addr)
		fields := logrus.Fields{"tx": tx.ID}
		if err != nil {
			a.log.WithFields(fields).Error("Buffer drain is failed")

			return errors.New("Buffer drain is failed")
		} else {
			a.log.WithFields(fields).Info("Buffer successfully drained")
		}
	} else {
		a.log.Debug("No need to drain buffer")
	}

	return nil
}

func (a *AdminActions) ChownStakeReDelegate(p *Protocol, showOnly bool) error {
	activeBuffer := p.GetActiveBuffer()
	aZilSsnAddr := p.Azil.GetAzilSsnAddress()

	a.log.WithFields(logrus.Fields{"active_buffer": activeBuffer.Addr}).Info("Active Buffer")

	mapSsnAmount := p.Zimpl.GetDepositAmtDeleg(activeBuffer.Addr)
	if 0 == len(mapSsnAmount) {
		a.log.Info("Buffer is empty, nothing to redelegate.")
		return nil
	}

	for ssn, amount := range mapSsnAmount {
		amountStr := amount.String()
		if ssn != aZilSsnAddr {
			if showOnly {
				a.log.WithFields(logrus.Fields{"ssn": ssn, "amount": amountStr}).Debug("Need to run ChownStakeReDelegate")
			} else if tx, err := p.Azil.ChownStakeReDelegate(ssn, amountStr); err != nil {
				a.log.WithFields(logrus.Fields{"ssn": ssn, "amount": amountStr, "tx": tx.ID}).Error("ChownStakeReDelegate failed")

				return errors.New("ChownStakeReDelegate failed")
			} else {
				a.log.WithFields(logrus.Fields{"ssn": ssn, "amount": amountStr, "tx": tx.ID}).Info("ChownStakeReDelegate OK")
			}
		} else {
			a.log.WithFields(logrus.Fields{"ssn": ssn, "amount": amountStr}).Info("Skip ChownStakeReDelegate")
		}
	}

	return nil
}

func (a *AdminActions) ConfirmSwapRequests(p *Protocol) error {
	nextBuffer := p.GetBufferToSwapWith().Addr
	swapRequests := p.GetSwapRequestsForBuffer(nextBuffer)
	a.log.WithFields(logrus.Fields{"swap_requests_count": len(swapRequests), "next_buffer": nextBuffer}).Debug("Swap request found")
	errCnt := 0
	okCnt := 0
	for _, initiator := range swapRequests {
		tx, err := p.Azil.ChownStakeConfirmSwap(initiator)
		if err != nil {
			a.log.WithFields(logrus.Fields{"initiator": initiator, "txid": tx.ID, "error": err.Error()}).Error("Can't confirm swap")
			errCnt++
		} else {
			a.log.WithFields(logrus.Fields{"initiator": initiator, "txid": tx.ID}).Info("Swap confirmed")
			okCnt++
		}
	}
	a.log.WithFields(logrus.Fields{"confirmed_swaps_count": okCnt, "errors_count": errCnt}).Debug("confirmSwapRequests completed")

	if errCnt > 0 {
		return errors.New("ConfirmSwapRequests has failed items")
	}

	return nil
}

func (a *AdminActions) AutoRestake(p *Protocol) error {
	autorestakeamount := p.Azil.GetAutorestakeAmount()

	if autorestakeamount.Cmp(big.NewInt(0)) == 0 { // autorestakeamount == 0
		a.log.Info("Nothing to autorestake")
		return nil
	}

	priceBefore := p.Azil.GetAzilPrice().String()
	tx, err := p.Azil.PerformAutoRestake()

	if err != nil {
		a.log.WithFields(logrus.Fields{"error": err.Error()}).Error("AutoRestake failed")

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
		a.log.WithFields(logrus.Fields{"tx": tx.ID}).Error("Withdrawals claim failed")

		return errors.New("Withdrawals claim failed")
	} else {
		a.log.WithFields(logrus.Fields{"tx": tx.ID}).Info("Withdrawals successfully claimed")
	}

	return nil
}
