package utils

import (
	"net/url"

	"github.com/Zilliqa/gozilliqa-sdk/subscription"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
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
	log *core.Log
}

func CreateBlockWatcher(url url.URL, log *core.Log) *BlockWatcher {
	return &BlockWatcher{url: url, log: log}
}

func (cw *BlockWatcher) Start() {
	subscriber := subscription.BuildNewBlockSubscriber(cw.url)
	_, ec, msg := subscriber.Start()
	cancel := false

	cw.log.WithFields(logrus.Fields{"url": cw.url.String()}).Debug("Start block watcher")

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
			cw.log.Error("Got error: " + err.Error())
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
