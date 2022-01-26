package main

import (
	"flag"
	"math/big"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/subscription"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/tidwall/gjson"
)

var log *Log
var sdk *AvelySDK

func tryDrainBuffer(p *Protocol, lrc int) {
	bufferToDrain := p.GetBufferToDrain()

	buffers := p.Aimpl.GetDrainedBuffers()
	needDrain := false

	if lastDrained, ok := buffers[strings.ToLower(bufferToDrain.Addr)]; ok {
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
			log.Fatal("Buffer drain is failed. Tx: " + tx.ID)
		} else {
			log.Success("Buffer successfully drained" + tx.ID)
		}
	} else {
		log.Success("No need to drain buffer")
	}
}

func tryAutorestake(p *Protocol) {
	autorestakeamount := p.Aimpl.GetAutorestakeAmount()

	if autorestakeamount.Cmp(big.NewInt(0)) == 0 { // autorestakeamount == 0
		log.Success("Nothing to autorestake")
	} else {
		priceBefore := p.Aimpl.GetAzilPrice()
		tx, err := p.Aimpl.PerformAutoRestake()

		if err != nil {
			log.Fatalf("AutoRestake failed with error: ", err)
		}

		priceAfter := p.Aimpl.GetAzilPrice()

		log.Success("AutoRestake is successfully completed. Tx: " + tx.ID)
		log.Success("Restaked amount: " + autorestakeamount.String() + "; PriceBefore: " + priceBefore.String() + "; PriceAfter: " + priceAfter.String())
	}
}

// If Last reward has been changed, then:
//   1. Drain Buffer
//   2. Autorestake funds
func main() {
	chainPtr := flag.String("chain", "local", "chain")

	flag.Parse()

	log = NewLog()
	config := NewConfig(*chainPtr)
	sdk = NewAvelySDK(*config)
	protocol := RestoreFromState(sdk, log)

	currentLrc := -1 // init value should be unreal

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
				blockNo := block.Get("BlockNum").String()
				lrc := protocol.GetLastRewardCycle()

				log.Success("New block #" + blockNo + " is mined. Current Last Reward Cycle is " + strconv.Itoa(lrc))

				if lrc == currentLrc {
					log.Success("Last Reward Cycle is not changed. Skip")
				} else {
					log.Success("New Last Reward Cycle!")
					tryDrainBuffer(protocol, lrc)
					tryAutorestake(protocol)

					currentLrc = lrc
				}
			}

		case err := <-ec:
			log.Success("Get error: ", err.Error())
			cancel = true
		}

	}
}
