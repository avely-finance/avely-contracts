package utils

import (
	"github.com/Zilliqa/gozilliqa-sdk/subscription"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/tidwall/gjson"
	"net/url"
)

type Observer interface {
	Notify(blockNum int)
}

type Observable struct {
	Observers []Observer
}

type BlockWatcher struct {
	Observable
	url url.URL
}

var log *core.Log

func CreateBlockWatcher(url url.URL) *BlockWatcher {
	log = core.NewLog()
	return &BlockWatcher{url: url}
}

func (cw *BlockWatcher) Start() {
	subscriber := subscription.BuildNewBlockSubscriber(cw.url)
	_, ec, msg := subscriber.Start()
	cancel := false

	log.Successf("Start block watcher, url=%s", cw.url.String())

	for {
		if cancel {
			_ = subscriber.Ws.Close()
			break
		}
		select {
		case message := <-msg:
			if blockNum, found := cw.getBlockNum(message); found {
				cw.NotifyAll(blockNum)
			}
		case err := <-ec:
			log.Errorf("Got error: %s", err)
			cancel = true
		}

	}
}

func (cw *BlockWatcher) getBlockNum(message []byte) (int, bool) {
	block := gjson.Get(string(message), "values.0.value.TxBlock.header")
	if block.Get("BlockNum").Exists() {
		blockNo := int(block.Get("BlockNum").Int())
		return blockNo, true
	}
	return -1, false
}

func (o *Observable) AddObserver(obs Observer) {
	o.Observers = append(o.Observers, obs)
}
func (o *Observable) NotifyAll(value int) {
	for _, ob := range o.Observers {
		ob.Notify(value)
	}
}