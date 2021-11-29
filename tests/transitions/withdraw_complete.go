package transitions

import (
	//"log"
	"Azil/test/deploy"
	"fmt"
	//"math/big"
)

func (t *Testing) CompleteWithdrawalSuccess() {

	t.LogStart("CompleteWithdrawal - success")

	stubStakingContract, aZilContract, _, _ := t.DeployAndUpgrade()
	t.AddDebug("addr1", "0x"+addr1)
	t.AddDebug("addr2", "0x"+addr2)

	aZilContract.UpdateWallet(key1)
	aZilContract.DelegateStake(zil10)
	stubStakingContract.AssignStakeReward()
	aZilContract.WithdrawStakeAmt(azil10)
	tx, _ := aZilContract.CompleteWithdrawal()
	t.AssertEvent(tx, aZilContract.Event("NoUnbondedStake", `{}`))

	aZilContract.UpdateWallet(key2)
	tx, _ = aZilContract.CompleteWithdrawal()
	t.AssertEvent(tx, aZilContract.Event("NoPendingWithdrawal", `{}`))

	deploy.IncreaseBlocknum(stubStakingContract.GetBnumReq() + 1)
	stubStakingContract.AssignStakeReward()

	aZilContract.UpdateWallet(key1)
	//balance1_before := deploy.GetBalance(addr1)
	tx, _ = aZilContract.CompleteWithdrawal()
	t.AssertEvent(tx, aZilContract.Event("CompleteWithdrawal", fmt.Sprintf(`{"amount": "%s"}`, zil10)))
	/*
		//https://github.com/Zilliqa/gozilliqa-sdk/blob/dd0ecada1be6987976b9f3b557dbb4de305ecf5b/account/wallet.go#L225
		balance1_after := deploy.GetBalance(addr1)
		after, _ := new(big.Int).SetString(balance1_after, 10)
		gasPrice, _ := new(big.Int).SetString(tx.GasPrice, 10)
		gasLimit, _ := new(big.Int).SetString(tx.GasLimit, 10)
		gasFee := new(big.Int).Mul(gasPrice, gasLimit)
		needed := new(big.Int).Add(gasFee, after)
		t.AssertEqual(balance1_before, fmt.Sprintf("%s", needed))*/

	/*
	 3) stub->changeBNUMreq(1), wait(5sec), completewithdr()
	     e = { _eventname: "CompleteWithdrawal"; amount: withdraw_amt }; проверить наличие, сравнить амаунт
	     msg = {_tag: "CompleteWithdrawal"; _recipient: holder_addr; _amount: uint128_zero };
	     msg = {_tag: "CompleteWithdrawal"; _recipient: proxy_staking_contract_addr; _amount: uint128_zero };
	*/

}
