package admin

import (
    . "github.com/avely-finance/avely-contracts/sdk/contracts"
    "github.com/avely-finance/avely-contracts/tests/helpers"
)

var log = helpers.GetLog()

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
