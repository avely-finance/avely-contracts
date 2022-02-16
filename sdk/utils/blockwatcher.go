package utils

import (
	"net/url"
	"time"

	"github.com/Zilliqa/gozilliqa-sdk/subscription"
	"github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const errElapse = float64(30.0)

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

type SocketError struct {
	err       error
	timestamp time.Time
}

func NewSocketError(e error) *SocketError {
	return &SocketError{err: e, timestamp: time.Now()}
}

func (s *SocketError) Error() string {
	return s.err.Error()
}

func (s *SocketError) isExpired() bool {
	now := time.Now()
	delta := now.Sub(s.timestamp)
	return delta.Seconds() > errElapse
}

func (cw *BlockWatcher) Start() {
	log := cw.log
	subscriber := subscription.BuildNewBlockSubscriber(cw.url)
	_, ec, msg := subscriber.Start()
	var prevErr *SocketError
	var err *SocketError

	log.WithFields(logrus.Fields{"url": cw.url.String()}).Debug("Start block watcher")

	for {
		if err != nil {
			_ = subscriber.Ws.Close()
			if prevErr == nil || prevErr.isExpired() {
				// reconnect
				subscriber = subscription.BuildNewBlockSubscriber(cw.url)
				_, ec, msg = subscriber.Start()
				// re-assign prev error to the current
				prevErr = err
				err = nil
			} else {
				log.Error("Got fatal error: " + err.Error())
				break
			}
		}
		select {
		case message := <-msg:
			if blockNum, found := cw.getBlockNum(message); found {
				cw.NotifyAll(blockNum)
			}
		case e := <-ec:
			err = NewSocketError(e)
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
