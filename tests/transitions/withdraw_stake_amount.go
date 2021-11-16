package transitions

import (
// "log"
    "regexp"
)

func (t *Testing) WithdrawStakeAmount() {
    t.LogStart("WithdrawStakeAmount")

    // deploy smart contract
    stubStakingContract, aZilContract, _, _ := t.DeployAndUpgrade()

    /*******************************************************************************
    * 0. delegator (addr2) delegate 10 zil, and it should enter in buffered deposit
    *******************************************************************************/
    aZilContract.UpdateWallet(key2)
    aZilContract.DelegateStake(fifteenzil)
    //we need to move buffered deposits to main stake
    stubStakingContract.AssignStakeReward()


    /*******************************************************************************
    * 1. non delegator(addr4) try to withdraw stake, should fail
    *******************************************************************************/
    t.LogStart("================== WithdwarStakeAmount, step 1 ===================" );
    aZilContract.UpdateWallet(key4)
    txn, err := aZilContract.WithdrawStakeAmt(tenzil)
    t.AssertError(err)
    t.LogPrettyReceipt(txn)
    receipt := t.GetReceiptString(txn)
    t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -7))])")


    /*******************************************************************************
    * 2A. delegator trying to withdraw more than staked, should fail
    *******************************************************************************/
    aZilContract.UpdateWallet(key2)
    t.LogStart("================== WithdwarStakeAmount, step 2A ===================" );
    txn, err = aZilContract.WithdrawStakeAmt(hundredzil)
    t.AssertError(err)
    t.LogPrettyReceipt(txn)
    receipt =  t.GetReceiptString(txn)
    t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -13))])")
    t.AssertContain(stubStakingContract.LogContractStateJson(), "\"totalstakeamount\":\"" + fifteenzil + "\"")


    /*******************************************************************************
    * 2B. delegator send withdraw request, but it should fail because mindelegatestake
    * TODO: how to be sure about size of mindelegatestake here?
    *******************************************************************************/
    t.LogStart("================== WithdwarStakeAmount, step 2B ===================" );
    txn, err = aZilContract.WithdrawStakeAmt(tenzil)
    t.AssertError(err)
    receipt = t.GetReceiptString(txn)
    t.AssertContain(receipt,"Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -15))])")
    t.AssertContain(stubStakingContract.LogContractStateJson(), "\"totalstakeamount\":\"" + fifteenzil + "\"")


    /*******************************************************************************
    * 3A. delegator withdrawing part of his deposit, it should success with "_eventname": "WithdrawStakeAmt"
    * Also check that withdrawal_pending field contains correct information about requested withdrawal
    * balances field should be correct
    *******************************************************************************/
    t.LogStart("================== WithdwarStakeAmount, step 3A ===================" );
    txn, err = aZilContract.WithdrawStakeAmt(fivezil)
    if err != nil {
        t.LogError("WithdrawStakeAmount",err)
    }
    t.AssertContain(t.GetReceiptString(txn), "WithdrawStakeAmt")
    aZilState := aZilContract.LogContractStateJson()
    t.AssertContain(aZilState, "\"totalstakeamount\":\"" + tenzil + "\",\"totaltokenamount\":\"" + tenzil + "\"")
    t.AssertContain(aZilState, "\"balances\":{\"" + "0x" + addr2 + "\":\"" + tenzil + "\"}");
    //replace epoch number with fake
    myRegexp := regexp.MustCompile(`\{\"(\d){1,10}\"\:\{\"argtypes\":\[\],`)
    aZilState = myRegexp.ReplaceAllString(aZilState, "{\"1234567890\":{\"argtypes\":[],")
    t.AssertContain(aZilState, "\"withdrawal_pending\":{\"" + "0x" + addr2 + "\":{\"" + /*txn.Receipt.EpochNum*/"1234567890" + "\":{\"argtypes\":[],\"arguments\":[\"" + fivezil + "\",\"" + fivezil + "\"]")
    t.AssertContain(stubStakingContract.LogContractStateJson(), "\"totalstakeamount\":\"" + tenzil + "\"")


    /*******************************************************************************
    * 3B. delegator withdrawing all remaining deposit, it should success with "_eventname": "WithdrawStakeAmt"
    * Also check that withdrawal_pending field contains correct information about requested withdrawal
    * Balances should be empty
    *******************************************************************************/
    t.LogStart("================== WithdrawStakeAmount, step 3B ===================" );
    txn, err = aZilContract.WithdrawStakeAmt(tenzil)
    if err != nil {
        t.LogError("WithdrawStakeAmount",err)
    }
    //check event
    t.AssertContain(t.GetReceiptString(txn), "WithdrawStakeAmt")
    t.AssertContain(t.GetReceiptString(txn), "{\"type\":\"Uint128\",\"value\":\"" + tenzil + "\",\"vname\":\"withdraw_amount\"},{\"type\":\"Uint128\",\"value\":\"" + tenzil + "\",\"vname\":\"withdraw_stake_amount\"}")
    //check contract state
    aZilState = aZilContract.LogContractStateJson()
    t.AssertContain(aZilState, "\"balances\":{},")
    t.AssertContain(aZilState, "\"totalstakeamount\":\"0\",\"totaltokenamount\":\"0\"")
    t.AssertContain(aZilState,"\"withdrawal_pending\":{\"" + "0x" + addr2 + "\":{\"" + txn.Receipt.EpochNum + "\":{\"argtypes\":[],\"arguments\":[\"" + fifteenzil + "\",\"" + fifteenzil + "\"]")
    t.AssertContain(stubStakingContract.LogContractStateJson(), "\"totalstakeamount\":\"" + "0" + "\"")

    // TODO: if delegator have buffered deposits, withdrawal should fail

    t.LogEnd("WithdrawStakeAmount")
}
