package actions

import (
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/avely-finance/avely-contracts/tests/helpers"
	"math/big"
	"strings"
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
		log.Success("Buffer is never drained; Let's do this first time")

		needDrain = true
	}

	if needDrain {
		tx, err := p.Aimpl.DrainBuffer(bufferToDrain.Addr)

		if err != nil {
			log.Fatalf("Buffer drain is failed. Tx: %s", tx.ID)
		} else {
			log.Successf("Buffer successfully drained. Tx: %s", tx.ID)
		}
	} else {
		log.Success("No need to drain buffer")
	}
}

func ChownStakeReDelegate(p *Protocol, showOnly bool) {
	activeBuffer := p.GetActiveBuffer()
	aZilSsnAddr := p.GetAzilSsnAddress()

	log.Successf("Active Buffer is %s", activeBuffer.Addr)

	mapSsnAmount := p.Zimpl.GetDepositAmtDeleg(activeBuffer.Addr)
	if 0 == len(mapSsnAmount) {
		log.Success("Buffer is empty, nothing to redelegate.")
		return
	}

	for ssn, amount := range mapSsnAmount {
		amountStr := amount.String()
		if ssn != aZilSsnAddr {
			if showOnly {
				log.Successf("SSN %s has %s. Need to run ChownStakeReDelegate.", ssn, amountStr)
			} else if tx, err := p.Aimpl.ChownStakeReDelegate(ssn, amountStr); err != nil {
				log.Fatalf("ChownStakeReDelegate(%s, %s) ERROR. Tx: %s.", ssn, amountStr, tx.ID)
			} else {
				log.Successf("ChownStakeReDelegate(%s, %s) OK, Tx: %s.", ssn, amountStr, tx.ID)
			}
		} else {
			log.Successf("AzilSSN %s has %s. Skip ChownStakeReDelegate.", ssn, amountStr)
		}
	}

}

func ConfirmSwapRequests(p *Protocol) {
	nextBuffer := p.GetBufferToSwapWith().Addr
	swapRequests := p.GetSwapRequestsForBuffer(nextBuffer)
	log.Successf("Found %d swap requests for next buffer %s", len(swapRequests), nextBuffer)
	errCnt := 0
	okCnt := 0
	for _, initiator := range swapRequests {
		tx, err := p.Aimpl.ChownStakeConfirmSwap(initiator)
		if err != nil {
			log.Errorf("Can't confirm swap: ChownStakeConfirmSwap(%s) error=%s, txid=%s", initiator, err, tx.ID)
			errCnt++
		} else {
			log.Successf("Swap confirmed: ChownStakeConfirmSwap(%s) OK, txid=%s", initiator, tx.ID)
			okCnt++
		}
	}
	log.Successf("confirmSwapRequests completed, %d swaps confirmed, %d errors", okCnt, errCnt)
}

func AutoRestake(p *Protocol) {
	autorestakeamount := p.Aimpl.GetAutorestakeAmount()

	if autorestakeamount.Cmp(big.NewInt(0)) == 0 { // autorestakeamount == 0
		log.Success("Nothing to autorestake")
		return
	}

	priceBefore := p.Aimpl.GetAzilPrice().String()
	tx, err := p.Aimpl.PerformAutoRestake()
	log.Info(tx)

	if err != nil {
		log.Fatalf("AutoRestake failed with error: %s", err)
	}

	priceAfter := p.Aimpl.GetAzilPrice().String()

	log.Successf("AutoRestake is successfully completed. Tx: %s.", tx.ID)
	log.Successf("Restaked amount: %s; PriceBefore: %s; PriceAfter: %s", autorestakeamount.String(), priceBefore, priceAfter)
}

func ShowClaimWithdrawal(p *Protocol) {
	blocks := p.GetClaimWithdrawalBlocks()
	if len(blocks) > 0 {
		blocksStr := utils.ArrayItoa(blocks)
		log.Successf("Blocks with unbonded withdrawals: %s.", strings.Join(blocksStr, ", "))
	} else {
		log.Successf("Blocks with unbonded withdrawals not found.")
	}
}

func ClaimWithdrawal(p *Protocol) {
	blocks := p.GetClaimWithdrawalBlocks()
	cnt := len(blocks)
	if cnt == 0 {
		log.Successf("There are no blocks with unbonded withdrawals.")
		return
	}
	blocksStr := utils.ArrayItoa(blocks)
	log.Successf("Found %d blocks with unbonded withdrawals: %s.", cnt, strings.Join(blocksStr, ", "))
	tx, err := p.Aimpl.ClaimWithdrawal(blocksStr)

	if err != nil {
		log.Fatalf("Withdrawals claim failed. Tx: %s", tx.ID)
	} else {
		log.Successf("Withdrawals successfully claimed. Tx: %s", tx.ID)
	}

}
