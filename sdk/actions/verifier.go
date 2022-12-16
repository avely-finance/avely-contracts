package actions

import (
	"github.com/avely-finance/avely-contracts/sdk/contracts"
	"github.com/avely-finance/avely-contracts/sdk/utils"
	"github.com/sirupsen/logrus"
)

func (a *AdminActions) NextCycleWithAmount(p *contracts.Protocol, amountPerSSN int) {
	totalAmount := 0

	zimplSsnList := p.Zimpl.GetSsnList()
	ssnRewardFactor := make(map[string]string)
	for _, ssn := range zimplSsnList {
		totalAmount += amountPerSSN
		ssnRewardFactor[ssn] = utils.ToQA(amountPerSSN)
	}

	cfg := p.Zproxy.Sdk.Cfg

	ssnRewardFactor[cfg.StZilSsnAddress] = cfg.StZilSsnRewardShare

	tx, err := p.Zproxy.AssignStakeRewardList(ssnRewardFactor, utils.ToQA(totalAmount))
	fields := logrus.Fields{
		"tx": tx.ID,
	}
	if err == nil {
		a.log.WithFields(fields).Info("Next Cycle trigger success")
	} else {
		fields["error"] = tx.Receipt
		a.log.WithFields(fields).Error("Next Cycle trigger error")
	}
}
