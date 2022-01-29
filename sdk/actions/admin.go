package actions

import (
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/tests/helpers"
)

var log = helpers.GetLog()

func ChownStakeReDelegate(p *Protocol) {
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
			tx, err := p.Aimpl.ChownStakeReDelegate(ssn, amountStr)
			if err != nil {
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
	log.Infof("Found %d swap requests for next buffer %s", len(swapRequests), nextBuffer)
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
	log.Infof("confirmSwapRequests completed, %d swaps confirmed, %d errors", okCnt, errCnt)
}
