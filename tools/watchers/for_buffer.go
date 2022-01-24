package main

import (
	"net/url"
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

	state := NewState(p.Aimpl.GetDrainedBuffers())
	buffers := state.Dig("result.balances").Map()
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

func main() {
	log = NewLog()
	config := NewConfig("testnet")
	sdk = NewAvelySDK(*config)
	protocol := RestoreFromState(sdk, log)

	currentLrc := -1 // default valus should be unreal

	u := url.URL{Scheme: "wss", Host: "dev-ws.zilliqa.com", Path: ""}
	subscriber := subscription.BuildNewBlockSubscriber(u)
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

				log.Success("New block #" + blockNo + " is mined. Current Last Reward Cycle is" + strconv.Itoa(lrc))

				if lrc == currentLrc {
					log.Success("Last Reward Cycle is not changed. Skip")
				} else {
					log.Success("New Last Reward Cycle!")
					tryDrainBuffer(protocol, lrc)

					currentLrc = lrc
				}
			}

		case err := <-ec:
			log.Success("Get error: ", err.Error())
			cancel = true
		}

	}
}
