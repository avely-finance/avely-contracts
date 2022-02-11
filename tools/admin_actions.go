package main

import (
	"flag"
	"fmt"
	"math/big"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/avely-finance/avely-contracts/sdk/actions"
	. "github.com/avely-finance/avely-contracts/sdk/contracts"
	. "github.com/avely-finance/avely-contracts/sdk/core"
	"github.com/sirupsen/logrus"
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
		//deploy
		case "init_holder":
			initHolder(p)
		case "sync_buffers":
			syncBuffers(p)
		case "deploy_buffer":
			deployBuffer(p)
		case "unpause":
			unpause(p)

		//utils, information
		case "from_bech32":
			convertFromBech32Addr(addr)
		case "to_bech32":
			convertToBech32Addr(addr)
		case "show_tx":
			showTx(p, addr)
		case "get_active_buffer":
			getActiveBuffer(p)
		case "show_rewards":
			showRewards(p, ssn, addr)

		//new reward cycle
		case "drain_buffer":
			drainBuffer(p, addr)
		case "redelegate":
			showOnly := false
			actions.ChownStakeReDelegate(p, showOnly)
		case "show_redelegate":
			showOnly := true
			actions.ChownStakeReDelegate(p, showOnly)
		case "autorestake":
			actions.AutoRestake(p)

		//withdrawals
		case "show_claim_withdrawal":
			actions.ShowClaimWithdrawal(p)
		case "claim_withdrawal":
			actions.ClaimWithdrawal(p)

		//swap requests (part of chown stake process)
		case "show_swap_requests":
			showSwapRequests(p)
		case "confirm_swap_requests":
			actions.ConfirmSwapRequests(p)

		default:
			log.Fatal("Unknown command")
		}
	}

	log.Info("Done")
}

func deployAvely() {
	p := DeployOnlyAvely(sdk, log)
	p.SyncBufferAndHolder()
}

func showTx(p *Protocol, tx_addr string) {
	provider := p.Aimpl.Contract.Provider
	tx, err := provider.GetTransaction(tx_addr)
	if err != nil {
		log.Error("Err: " + err.Error())
	} else {
		log.Info(tx)
	}
}

func getActiveBuffer(p *Protocol) {
	log.WithFields(logrus.Fields{
		"lrc":           p.Zimpl.GetLastRewardCycle(),
		"active_buffer": p.GetActiveBuffer().Contract.Addr,
	}).Info("Active buffer / lrc")
}

func initHolder(p *Protocol) {
	_, err := p.InitHolder()

	if err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Fatal("Holder init failed")
	}
	log.Info("Holder init is successfully compeleted.")
}

func convertFromBech32Addr(addr32 string) {
	addr, err := bech32.FromBech32Addr(addr32)

	if err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Fatal("Convert failed")
	}

	log.WithFields(logrus.Fields{"addr": addr}).Info("Converted address")
}

func convertToBech32Addr(addr32 string) {
	addr, err := bech32.ToBech32Address(addr32)

	if err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Fatal("Convert failed")
	}

	log.WithFields(logrus.Fields{"addr": addr}).Info("Converted address")
}

func deployBuffer(p *Protocol) {
	buffer, err := p.DeployBuffer()

	if err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Fatal("Buffer deploy failed")
	}
	log.WithFields(logrus.Fields{"address": buffer.Addr}).Info("Buffer deploy is successfully completed")
}

func unpause(p *Protocol) {
	_, err := p.Aimpl.Unpause()

	if err != nil {
		log.WithFields(logrus.Fields{"error": err.Error()}).Fatal("Unpause AZil failed")
	}
	log.Info("Unpause AZil is successfully completed")
}

func syncBuffers(p *Protocol) {
	p.SyncBufferAndHolder()
}

func drainBuffer(p *Protocol, buffer_addr string) {
	tx, err := p.Aimpl.DrainBuffer(buffer_addr)

	if err != nil {
		log.WithFields(logrus.Fields{"buffer_addr": buffer_addr, "error": err.Error()}).Fatal("Drain buffer failed")
	}
	log.WithFields(logrus.Fields{"buffer_addr": buffer_addr, "tx": tx.ID}).Info("Drain is successfully completed")
}

func showRewards(p *Protocol, ssn, deleg string) {
	// result := p.Aimpl.Contract.SubState("balances",  [1]string{"0x79c7e38dd3b3c88a3fb182f26b66d8889e61cbd6"})

	rawState := p.Zimpl.Contract.State()

	state := NewState(rawState)

	one := big.NewInt(1)
	lastWithdrawCycle := state.Dig("last_withdraw_cycle_deleg", deleg, ssn).BigInt()
	lrc := big.NewInt(int64(p.Zimpl.GetLastRewardCycle()))

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

func showSwapRequests(p *Protocol) {
	i := 0
	partialState := p.Zimpl.Contract.SubState("deleg_swap_request", []string{})
	state := NewState(partialState)
	swapRequests := state.Dig("result.deleg_swap_request").Map()
	nextBuffer := p.GetBufferToSwapWith().Addr
	buffers := make(map[string]bool)
	for _, buffer := range p.Buffers {
		buffers[buffer.Addr] = true
	}
	log.WithFields(logrus.Fields{"active_buffer": p.GetActiveBuffer().Addr, "next_buffer": nextBuffer}).Info("Buffers")
	for initiator, new_deleg := range swapRequests {
		new_deleg_string := new_deleg.String()
		status := ""
		_, isBuffer := buffers[new_deleg_string]
		if isBuffer {
			i++
			switch new_deleg_string {
			case nextBuffer:
				status = "OK"
			default:
				status = "WRONG BUFFER"
			}
			log.WithFields(logrus.Fields{"initiator": initiator, "new_deleg": new_deleg_string, "status": status}).Info("Swap request")
		}
	}
	log.WithFields(logrus.Fields{"count": i}).Info("Swap requests total")
}
