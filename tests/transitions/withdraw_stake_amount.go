package transitions

import (
    // "log"
    "regexp"
)

func (t *Testing) WithdrawStakeAmount() {

    const FakeEpochNum = "1234567890"

    t.LogStart("WithdrawStakeAmount")

    // deploy smart contract
    stubStakingContract, aZilContract, _, _ := t.DeployAndUpgrade()

    /*******************************************************************************
     * 0. delegator (addr2) delegate 15 zil, and it should enter in buffered deposit,
     * we need to move buffered deposits to main stake
     *******************************************************************************/
    aZilContract.UpdateWallet(key2)
    aZilContract.DelegateStake(unit15)
    // TODO: if delegator have buffered deposits, withdrawal should fail
    stubStakingContract.AssignStakeReward()

    /*******************************************************************************
     * 1. non delegator(addr4) try to withdraw stake, should fail
     *******************************************************************************/
    t.LogStart("================== WithdwarStakeAmount, step 1 ===================")
    aZilContract.UpdateWallet(key4)
    txn, err := aZilContract.WithdrawStakeAmt(unit10)
    t.AssertError(err)
    t.LogPrettyReceipt(txn)
    t.AssertContain(t.GetReceiptString(txn), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -7))])")

    /*******************************************************************************
     * 2A. delegator trying to withdraw more than staked, should fail
     *******************************************************************************/
    aZilContract.UpdateWallet(key2)
    t.LogStart("================== WithdwarStakeAmount, step 2A ===================")
    txn, err = aZilContract.WithdrawStakeAmt(unit100)
    t.AssertError(err)
    t.LogPrettyReceipt(txn)
    t.AssertContain(t.GetReceiptString(txn), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -13))])")
    t.AssertContain(aZilContract.LogContractStateJson(), "\"totaltokenamount\":\""+unit15+"\"")

    /*******************************************************************************
     * 2B. delegator send withdraw request, but it should fail because mindelegatestake
     * TODO: how to be sure about size of mindelegatestake here?
     *******************************************************************************/
    t.LogStart("================== WithdwarStakeAmount, step 2B ===================")
    txn, err = aZilContract.WithdrawStakeAmt(unit10)
    t.AssertError(err)
    t.AssertContain(t.GetReceiptString(txn), "Exception thrown: (Message [(_exception : (String \\\"Error\\\")) ; (code : (Int32 -15))])")
    t.AssertContain(aZilContract.LogContractStateJson(), "\"totaltokenamount\":\""+unit15+"\"")

    /*******************************************************************************
     * 3A. delegator withdrawing part of his deposit, it should success with "_eventname": "WithdrawStakeAmt"
     * Also check that withdrawal_pending field contains correct information about requested withdrawal
     * balances field should be correct
     *******************************************************************************/
    t.LogStart("================== WithdwarStakeAmount, step 3A ===================")
    txn, err = aZilContract.WithdrawStakeAmt(unit5)
    if err != nil {
        t.LogError("WithdrawStakeAmount", err)
    }
    t.AssertContain(t.GetReceiptString(txn), "WithdrawStakeAmt")
    newDelegBalanceZil, err := aZilContract.ZilBalanceOf(addr2)
    aZilState := aZilContract.LogContractStateJson()
    if nil != err {
        t.LogError("WithdrawStakeAmount", err)
    }
    t.AssertContain(aZilState, "\"totalstakeamount\":\""+newDelegBalanceZil+"\",\"totaltokenamount\":\""+unit10+"\"")
    t.AssertContain(aZilState, "\"balances\":{\""+"0x"+addr2+"\":\""+unit10+"\"}")
    //replace epoch number with fake
    myRegexp := regexp.MustCompile(`\{\"(\d){1,10}\"\:\{\"argtypes\":\[\],`)
    aZilState = myRegexp.ReplaceAllString(aZilState, "{\""+FakeEpochNum+"\":{\"argtypes\":[],")
    t.AssertContain(aZilState, "\"withdrawal_pending\":{\""+"0x"+addr2+"\":{\""+ /*txn.Receipt.EpochNum*/ FakeEpochNum+"\":{\"argtypes\":[],\"arguments\":[\""+unit5+"\",\""+unit5+"\"]")
    t.AssertContain(stubStakingContract.LogContractStateJson(), "\"totalstakeamount\":\""+newDelegBalanceZil+"\"")

    /*******************************************************************************
     * 3B. delegator withdrawing all remaining deposit, it should success with "_eventname": "WithdrawStakeAmt"
     * Also check that withdrawal_pending field contains correct information about requested withdrawal
     * Balances should be empty
     *******************************************************************************/
    t.LogStart("================== WithdrawStakeAmount, step 3B ===================")
    txn, err = aZilContract.WithdrawStakeAmt(unit10)
    if err != nil {
        t.LogError("WithdrawStakeAmount", err)
    }
    //check event
    t.AssertContain(t.GetReceiptString(txn), "WithdrawStakeAmt")
    t.AssertContain(t.GetReceiptString(txn), "{\"type\":\"Uint128\",\"value\":\""+unit10+"\",\"vname\":\"withdraw_amount\"},{\"type\":\"Uint128\",\"value\":\""+unit10+"\",\"vname\":\"withdraw_stake_amount\"}")
    //check contract state
    aZilState = aZilContract.LogContractStateJson()
    t.AssertContain(aZilState, "\"balances\":{},")
    t.AssertContain(aZilState, "\"totalstakeamount\":\"0\",\"totaltokenamount\":\"0\"")
    t.AssertContain(stubStakingContract.LogContractStateJson(), "\"totalstakeamount\":\""+"0"+"\"")
    /* this assertion is commented, because subsequent withdrawals may go to different block, so it's not trivial to check total withdrawals amount
       * seems it's enough that we check withdrawal_pending at previous tests and zero-total here
       //replace epoch number with fake
       myRegexp = regexp.MustCompile(`\{\"(\d){1,10}\"\:\{\"argtypes\":\[\],`)
       aZilState = myRegexp.ReplaceAllString(aZilState, "{\"" + FakeEpochNum + "\":{\"argtypes\":[],")
       t.AssertContain(aZilState,"\"withdrawal_pending\":{\"" + "0x" + addr2 + "\":{\"" + FakeEpochNum + "\":{\"argtypes\":[],\"arguments\":[\"" + unit15 + "\",\"" + unit15 + "\"]")
    */

    t.LogEnd("WithdrawStakeAmount")
}
