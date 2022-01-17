package main

import (
	"flag"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"math/big"
	"strings"
)

var log *Log
var sdk *AvelySDK

func main() {
	chainPtr := flag.String("chain", "local", "chain")
	cmdPtr := flag.String("cmd", "default", "specific command")
	addrPtr := flag.String("addr", "default", "an entity address")
	ssnPtr := flag.String("ssn", "default", "an entity ssn address")

	flag.Parse()

	cmd := *cmdPtr

	log = NewLog()
	config := NewConfig(*chainPtr)
	sdk = NewAvelySDK(*config)

	shortcuts := map[string]string{
		"azilssn":  config.AzilSsnAddress,
		"addr1":    config.Addr1,
		"addr2":    config.Addr2,
		"addr3":    config.Addr3,
		"admin":    config.Admin,
		"verifier": config.Verifier,
	}
	log.AddShortcuts(shortcuts)

	if cmd == "deploy" {
		deployAvely()
	} else {
		// for non-deploy commands we need initialize protocol from config
		p := RestoreFromState(sdk, log)
		addr := strings.ToLower(*addrPtr)
		ssn := strings.ToLower(*ssnPtr)

		switch cmd {
		case "from_bech32":
			convertFromBech32Addr(addr)
		case "to_bech32":
			convertToBech32Addr(addr)
		case "show_tx":
			showTx(p, addr)
		case "get_active_buffer":
			getActiveBuffer(p)
		case "init_holder":
			initHolder(p)
		case "deploy_buffer":
			deployBuffer(p)
		case "unpause":
			unpause(p)
		case "sync_buffers":
			syncBuffers(p)
		case "drain_buffer":
			drainBuffer(p, addr)
		case "show_rewards":
			showRewards(p, ssn, addr)
		case "show_swap_requests":
			showSwapRequests(p)
		case "autorestake":
			autorestake(p)
		default:
			log.Fatal("Unknown command")
		}
	}

	log.Success("Done")
}

func deployAvely() {
	p := DeployOnlyAvely(sdk, log)
	p.SyncBufferAndHolder()
}

func showTx(p *Protocol, tx_addr string) {
	provider := p.Aimpl.Contract.Provider
	tx, err := provider.GetTransaction(tx_addr)

	log.Successf("Tx: ", tx)
	log.Successf("Err: ", err)
}

func getActiveBuffer(p *Protocol) {
	lrc, buffer := p.GetActiveBuffer()

	log.Successf("lastrewardcycle: ", lrc)
	log.Successf("Active buffer: ", buffer.Contract.Addr)
}

func initHolder(p *Protocol) {
	_, err := p.InitHolder()

	if err != nil {
		log.Fatalf("Holder init failed with error: ", err)
	}
	log.Success("Holder init is successfully compeleted.")
}

func convertFromBech32Addr(addr32 string) {
	addr, err := bech32.FromBech32Addr(addr32)

	if err != nil {
		log.Fatalf("Convert failed with err: ", err)
	}

	log.Success("Converted address: " + addr)
}

func convertToBech32Addr(addr32 string) {
	addr, err := bech32.ToBech32Address(addr32)

	if err != nil {
		log.Fatalf("Convert failed with err: ", err)
	}

	log.Success("Converted address: " + addr)
}

func deployBuffer(p *Protocol) {
	buffer, err := p.DeployBuffer()

	if err != nil {
		log.Fatalf("Buffer deploy failed with error: ", err)
	}
	log.Success("Buffer deploy is successfully completed. Address: " + buffer.Addr)
}

func unpause(p *Protocol) {
	_, err := p.Aimpl.Unpause()

	if err != nil {
		log.Fatalf("Unpause AZil failed with error: ", err)
	}
	log.Success("Unpause AZil is successfully completed")
}

func syncBuffers(p *Protocol) {
	p.SyncBufferAndHolder()
}

func drainBuffer(p *Protocol, buffer_addr string) {
	tx, err := p.Aimpl.DrainBuffer(buffer_addr)

	if err != nil {
		log.Fatalf("Drain failed with error: ", err)
	}
	log.Success("Drain is successfully completed. Tx: " + tx.ID)
}

func showRewards(p *Protocol, ssn, deleg string) {
	// result := p.Aimpl.Contract.SubState("balances",  [1]string{"0x79c7e38dd3b3c88a3fb182f26b66d8889e61cbd6"})

	rawState := p.Zimpl.Contract.State()

	state := NewState(rawState)

	one := big.NewInt(1)
	lastWithdrawCycle := state.Dig("last_withdraw_cycle_deleg", deleg, ssn).BigInt()
	lrc := state.Dig("lastrewardcycle").BigInt()

	m := AddBI(lastWithdrawCycle, one) // + 1
	n := lrc                           // iota should not include the last cycle since it does not completed yet

	delegStakePerCycle := big.NewInt(0)
	cycleRewardsDeleg := big.NewInt(0)

	// for cycle := m; cycle < n; cycle++
	for cycle := new(big.Int).Set(m); cycle.Cmp(n) < 0; cycle.Add(cycle, one) {

		// last_reward_cycle = builtin sub reward_cycle uint32_one;
		lastRewardCycle := SubBI(cycle, one)

		// last2_reward_cycle = sub_one_to_zero last_reward_cycle;
		last2RewardCycle := SubOneToZero(lastRewardCycle)

		// cur_opt <- direct_deposit_deleg[deleg][ssn_operator][last_reward_cycle];
		curOpt := state.Dig("direct_deposit_deleg", deleg, ssn, lastRewardCycle.String()).BigInt()
		// buf_opt <- buff_deposit_deleg[deleg][ssn_operator][last2_reward_cycle];
		bufOpt := state.Dig("buff_deposit_deleg", deleg, ssn, last2RewardCycle.String()).BigInt()
		// comb_opt = option_add cur_opt buf_opt;
		combOpt := AddBI(curOpt, bufOpt)

		// staking_of_deleg = match comb_opt with
		// | Some stake => builtin add last_amt stake
		// | None => last_amt
		// end;
		delegStakePerCycle = AddBI(delegStakePerCycle, combOpt)

		// staking_and_rewards_per_cycle_for_ssn_opt <- stake_ssn_per_cycle[ssn_operator][reward_cycle];
		rewardsForSsn := state.Dig("stake_ssn_per_cycle", ssn, cycle.String()).SSNCycleInfo()

		// fmt.Println(rewardsForSsn)
		reward := big.NewInt(0)

		if rewardsForSsn.TotalStaking.Cmp(big.NewInt(0)) == 1 { // TotalStaking > 0
			// reward = muldiv total_rewards staking_of_deleg total_staking;
			reward = reward.Mul(rewardsForSsn.TotalRewards, delegStakePerCycle)
			reward.Div(reward, rewardsForSsn.TotalStaking)
		} else {
			fmt.Println("SSN Node was excluded from this reward cycle")
		}

		cycleRewardsDeleg = AddBI(cycleRewardsDeleg, reward)
		fmt.Println("The cycle is: " + cycle.String() + "; Total reward: " + cycleRewardsDeleg.String())
		fmt.Println("    cur_opt: " + curOpt.String() + "; buf_opt: " + bufOpt.String())
	}
}

func autorestake(p *Protocol) {
	state := NewState(p.Aimpl.Contract.State())

	autorestakeamount := state.Dig("autorestakeamount").BigInt()

	if autorestakeamount.Cmp(big.NewInt(0)) == 0 { // == 0
		log.Fatal("Nothing to auto restake")
	}

	totaltokenamount := state.Dig("totaltokenamount").BigFloat()
	totalstakeamount := state.Dig("totalstakeamount").BigFloat()

	priceBefore := DivBF(totalstakeamount, totaltokenamount)

	tx, err := p.Aimpl.PerformAutoRestake()

	if err != nil {
		log.Fatalf("AutoRestake failed with error: ", err)
	}

	state = NewState(p.Aimpl.Contract.State())

	totalstakeamount = state.Dig("totalstakeamount").BigFloat()

	priceAfter := DivBF(totalstakeamount, totaltokenamount)

	log.Success("Drain is successfully completed. Tx: " + tx.ID)
	log.Success("Restaked amount: " + autorestakeamount.String() + "; PriceBefore: " + priceBefore.String() + "; PriceAfter: " + priceAfter.String())
}

func showSwapRequests(p *Protocol) {
	state := NewState(p.Zimpl.Contract.State())
	swapRequests := state.Dig("deleg_swap_request").Map()
	i := 0
	for initiator, new_deleg := range swapRequests {
		if new_deleg.String() == p.Holder.Addr {
			i++
			log.Info(initiator)
		}
		log.Infof("Found %d swap request(s)", i)
	}
}
