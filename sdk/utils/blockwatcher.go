package utils

import (
	"time"

	"github.com/avely-finance/avely-contracts/sdk/core"
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
	sdk *core.AvelySDK
	log *core.Log
}

func CreateBlockWatcher(sdk *core.AvelySDK, log *core.Log) *BlockWatcher {
	return &BlockWatcher{sdk: sdk, log: log}
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

	log.Debug("Block short polling watcher is ready to tick every 10 seconds")
	ticker := time.NewTicker(time.Second * 10).C

	go func() {
		for {
			select {
			case <-ticker:
				blockNum, err := cw.sdk.GetBlockHeight()

				if err != nil {
					log.Error("Got fatal error during fetchin block height: " + err.Error())
				}

				cw.NotifyAll(blockNum)
			}
		}
	}()

	time.Sleep(time.Duration(1<<63 - 1))
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
