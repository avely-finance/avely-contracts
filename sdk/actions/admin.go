package actions

import (
	"math/big"
	"strings"

	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/avely-finance/avely-contracts/tests/helpers"
	"github.com/sirupsen/logrus"
)

var log = helpers.GetLog()

func DrainBuffer(p *Protocol, lrc int) {
	bufferToDrain := p.GetBufferToDrain()

	buffers := p.Aimpl.GetDrainedBuffers()
	needDrain := false

	if lastDrained, ok := buffers[bufferToDrain.Addr]; ok {
		if lastDrained.Int() != int64(lrc) {
			needDrain = true
		}
	} else {
		log.Info("Buffer is never drained; Let's do this first time")

		needDrain = true
	}

	if needDrain {
		tx, err := p.Aimpl.DrainBuffer(bufferToDrain.Addr)
		fields := logrus.Fields{"tx": tx.ID}
		if err != nil {
			log.WithFields(fields).Fatal("Buffer drain is failed")
		} else {
			log.WithFields(fields).Info("Buffer successfully drained")
		}
	} else {
		log.Info("No need to drain buffer")
	}
}

func ChownStakeReDelegate(p *Protocol, showOnly bool) {
	activeBuffer := p.GetActiveBuffer()
	aZilSsnAddr := p.Aimpl.GetAzilSsnAddress()

	log.WithFields(logrus.Fields{"active_buffer": activeBuffer.Addr}).Info("Active Buffer")

	mapSsnAmount := p.Zimpl.GetDepositAmtDeleg(activeBuffer.Addr)
	if 0 == len(mapSsnAmount) {
		log.Info("Buffer is empty, nothing to redelegate.")
		return
	}

	for ssn, amount := range mapSsnAmount {
		amountStr := amount.String()
		if ssn != aZilSsnAddr {
			if showOnly {
				log.WithFields(logrus.Fields{"ssn": ssn, "amount": amountStr}).Debug("Need to run ChownStakeReDelegate")
			} else if tx, err := p.Aimpl.ChownStakeReDelegate(ssn, amountStr); err != nil {
				log.WithFields(logrus.Fields{"ssn": ssn, "amount": amountStr, "tx": tx.ID}).Fatal("ChownStakeReDelegate Error")
			} else {
				log.WithFields(logrus.Fields{"ssn": ssn, "amount": amountStr, "tx": tx.ID}).Info("ChownStakeReDelegate OK")
			}
		} else {
			log.WithFields(logrus.Fields{"ssn": ssn, "amount": amountStr}).Info("Skip ChownStakeReDelegate")
		}
	}

}

func ConfirmSwapRequests(p *Protocol) {
	nextBuffer := p.GetBufferToSwapWith().Addr
	swapRequests := p.GetSwapRequestsForBuffer(nextBuffer)
	log.WithFields(logrus.Fields{"swap_requests_count": len(swapRequests), "next_buffer": nextBuffer}).Debug("Swap request found")
	errCnt := 0
	okCnt := 0
	for _, initiator := range swapRequests {
		tx, err := p.Aimpl.ChownStakeConfirmSwap(initiator)
		if err != nil {
			log.WithFields(logrus.Fields{"initiator": initiator, "txid": tx.ID, "error": err.Error()}).Error("Can't confirm swap")
			errCnt++
		} else {
			log.WithFields(logrus.Fields{"initiator": initiator, "txid": tx.ID}).Info("Swap confirmed")
			okCnt++
		}
	}
	log.WithFields(logrus.Fields{"confirmed_swaps_count": okCnt, "errors_count": errCnt}).Debug("confirmSwapRequests completed")
}

func AutoRestake(p *Protocol) {
	autorestakeamount := p.Aimpl.GetAutorestakeAmount()

	if autorestakeamount.Cmp(big.NewInt(0)) == 0 { // autorestakeamount == 0
		log.Info("Nothing to autorestake")
		return
	}

	priceBefore := p.Aimpl.GetAzilPrice().String()
	tx, err := p.Aimpl.PerformAutoRestake()

	if err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Fatal("AutoRestake failed")
	}

	priceAfter := p.Aimpl.GetAzilPrice().String()

	log.WithFields(logrus.Fields{
		"tx":           tx.ID,
		"amount":       autorestakeamount.String(),
		"price_before": priceBefore,
		"price_after":  priceAfter,
	}).Info("AutoRestake is successfully completed")
}

func ShowClaimWithdrawal(p *Protocol) {
	blocks := p.GetClaimWithdrawalBlocks()
	if len(blocks) > 0 {
		blocksStr := utils.ArrayItoa(blocks)
		log.Debug("Blocks with unbonded withdrawals: " + strings.Join(blocksStr, ", "))
	} else {
		log.Debug("Blocks with unbonded withdrawals not found.")
	}
}

func ClaimWithdrawal(p *Protocol) {
	blocks := p.GetClaimWithdrawalBlocks()
	cnt := len(blocks)
	if cnt == 0 {
		log.Debug("There are no blocks with unbonded withdrawals.")
		return
	}
	blocksStr := utils.ArrayItoa(blocks)
	log.WithFields(logrus.Fields{"blocks_count": cnt, "blocks_list": strings.Join(blocksStr, ", ")}).Debug("Found blocks with unbonded withdrawals")
	tx, err := p.Aimpl.ClaimWithdrawal(blocksStr)

	if err != nil {
		log.WithFields(logrus.Fields{"tx": tx.ID}).Fatal("Withdrawals claim failed")
	} else {
		log.WithFields(logrus.Fields{"tx": tx.ID}).Info("Withdrawals successfully claimed")
	}

}
