package transitions

import (
	//"log"
	"Azil/test/deploy"
	//"math/big"
)

func (t *Testing) CompleteWithdrawalSuccess() {

	t.LogStart("CompleteWithdrawal - success")

	stubStakingContract, aZilContract, _, holderContract := t.DeployAndUpgrade()
	t.AddDebug("addr1", "0x"+addr1)
	t.AddDebug("addr2", "0x"+addr2)

	aZilContract.UpdateWallet(key1)
	aZilContract.DelegateStake(zil10)
	stubStakingContract.AssignStakeReward()
	aZilContract.WithdrawStakeAmt(azil10)
	tx, _ := aZilContract.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{aZilContract.Addr, "NoUnbondedStake", deploy.ParamsMap{}})

	aZilContract.UpdateWallet(key2)
	tx, _ = aZilContract.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{aZilContract.Addr, "NoPendingWithdrawal", deploy.ParamsMap{}})

	deploy.IncreaseBlocknum(stubStakingContract.GetBnumReq() + 1)
	stubStakingContract.AssignStakeReward()

	aZilContract.UpdateWallet(key1)
	//balance1_before := deploy.GetBalance(addr1)
	tx, _ = aZilContract.CompleteWithdrawal()
	t.AssertEvent(tx, deploy.Event{aZilContract.Addr, "CompleteWithdrawal", deploy.ParamsMap{"amount": zil10}})

	t.AssertTransition(tx, deploy.Transition{
		aZilContract.Addr,    //sender
		"CompleteWithdrawal", //tag
		holderContract.Addr,  //recipient
		"0",                  //amount
		deploy.ParamsMap{},
	})
	t.AssertTransition(tx, deploy.Transition{
		aZilContract.Addr,
		"CompleteWithdrawalSuccessCallBack",
		addr1,
		"0",
		deploy.ParamsMap{"amount": zil10},
	})

	/*
		before, _ := new(big.Int).SetString(balance1_before, 10)
		after, _ := new(big.Int).SetString(deploy.GetBalance(addr1), 10)
		gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
		gasUsed, _ := new(big.Int).SetString(tx.Receipt.CumulativeGas, 10)
		gasFee := new(big.Int).Mul(gasPrice, gasUsed)
		withdrawed, _ := new(big.Int).SetString(zil10, 10)
		result := new(big.Int).Sub(before, gasFee)
		result = result.Add(result, withdrawed)
		result = result.Sub(result, after)
		println(fmt.Sprintf("====%s====", result))
		t.LogDebug()
		//t.AssertEqual(balance1_before, fmt.Sprintf("%s", needed))*/

	/*
	 3) stub->changeBNUMreq(1), wait(5sec), completewithdr()
	     e = { _eventname: "CompleteWithdrawal"; amount: withdraw_amt }; проверить наличие, сравнить амаунт
	     msg = {_tag: "CompleteWithdrawal"; _recipient: holder_addr; _amount: uint128_zero };
	     msg = {_tag: "CompleteWithdrawal"; _recipient: proxy_staking_contract_addr; _amount: uint128_zero };
	*/

}
