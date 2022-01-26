package main

import (
	"flag"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/subscription"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/tidwall/gjson"
)

var log *Log
var sdk *AvelySDK

func tryRedelegateStakes(p *Protocol) {
	activeBuffer := p.GetActiveBuffer()
	aZilSsnAddr := p.GetAzilSsnAddress()

	log.Success("Active Buffer is " + activeBuffer.Addr)
	addr := strings.ToLower(activeBuffer.Addr)
	deposits, _ := p.Zimpl.GetDeposiAmdDeleg(addr)

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

// Every N blocks watch Redelegate Stakes
func main() {
	chainPtr := flag.String("chain", "local", "chain")
	gapPtr := flag.Int("gap", 5, "gap between blocks")

	flag.Parse()

	log = NewLog()
	config := NewConfig(*chainPtr)
	sdk = NewAvelySDK(*config)
	protocol := RestoreFromState(sdk, log)

	currentBlock := 0 // init value should be unreal

	url := sdk.GetWsURL()
	subscriber := subscription.BuildNewBlockSubscriber(url)
	_, ec, msg := subscriber.Start()
	cancel := false

	log.Success("Start subscribsion")

	for {
		if cancel {
			_ = subscriber.Ws.Close()
			break
		}
		select {
		case message := <-msg:
			block := gjson.Get(string(message), "values.0.value.TxBlock.header")

			if block.Get("BlockNum").Exists() {
				blockNo := int(block.Get("BlockNum").Int())

				log.Success("New block #" + strconv.Itoa(blockNo) + " is mined.")

				if (blockNo - currentBlock) > *gapPtr {
					tryRedelegateStakes(protocol)

					currentBlock = blockNo
				} else {
					log.Success("Gap is too low. Skip")
				}
			}

		case err := <-ec:
			log.Success("Get error: ", err.Error())
			cancel = true
		}

	}
}
