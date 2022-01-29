package actions

import (
    . "github.com/avely-finance/avely-contracts/sdk/contracts"
    "github.com/avely-finance/avely-contracts/tests/helpers"
    "strings"
)

var log = helpers.GetLog()

func ChownStakeReDelegate(p *Protocol) {
    activeBuffer := p.GetActiveBuffer()
    aZilSsnAddr := p.GetAzilSsnAddress()

    log.Success("Active Buffer is " + activeBuffer.Addr)
    addr := strings.ToLower(activeBuffer.Addr)
    deposits, _ := p.Zimpl.GetDeposiAmtDeleg(addr)

    if value, found := deposits[addr]; found {
        depositsMap := value.Map()
        for ssn, amount := range depositsMap {
            if ssn != aZilSsnAddr {
                tx, err := p.Aimpl.ChownStakeReDelegate(ssn, amount.String())

                if err != nil {
                    log.Fatal("ChownStakeReDelegate is failed. Tx: " + tx.ID)
                } else {
                    log.Success("Successfully redelegate " + amount.String() + " from SSN " + ssn + "; Tx: " + tx.ID)
                }
            }

            log.Success(ssn + " has " + amount.String())
        }
    } else {
        log.Success("Buffer is empty")
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
