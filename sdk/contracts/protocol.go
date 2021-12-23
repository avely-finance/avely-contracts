package contracts

import(
	"fmt"
	"log"
	"strconv"
	"runtime"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/core"
)

type Protocol struct {
	Zproxy *Zproxy
	Zimpl *Zimpl
	Aimpl *AZil
	Buffer *BufferContract
	Holder *HolderContract
}

func NewProtocol(zproxy *Zproxy, zimpl *Zimpl, azil *AZil, buffer *BufferContract, holder *HolderContract) *Protocol {
	return &Protocol{
		Zproxy: zproxy,
		Zimpl: zimpl,
		Aimpl: azil,
		Buffer: buffer,
		Holder: holder,
	}
}

func (p *Protocol) SyncBufferAndHolder() {
  new_buffers := []string{"0x" + p.Buffer.Addr}

	check(p.Aimpl.ChangeBuffers(new_buffers))
	check(p.Aimpl.ChangeHolderAddress(p.Holder.Addr))
}

func (p *Protocol) SetupZProxy() {
	sdk := p.Aimpl.Sdk
	args := []core.ContractValue{
		{
			"newImplementation",
			"ByStr20",
			"0x" + p.Zimpl.Addr,
		},
	}
	check(p.Zproxy.Call("UpgradeTo", args, "0"))
	check(p.Zproxy.AddSSN(sdk.Cfg.AzilSsnAddress, "aZil SSN"))
	check(p.Zproxy.UpdateVerifierRewardAddr("0x" + sdk.Cfg.Verifier))
	check(p.Zproxy.UpdateVerifier("0x" + sdk.Cfg.Verifier))
	check(p.Zproxy.UpdateStakingParameters(zil(1000), zil(10))) //minstake (ssn not active if less), mindelegstake
	check(p.Zproxy.Unpause())

	//we need our SSN to be active, so delegating some stake
	check(p.Aimpl.DelegateStake(zil(1000)))

	//we need to delegate something from Holder, in order to make Zimpl know holder's address
	check(p.Holder.DelegateStake(zil(sdk.Cfg.HolderInitialDelegateZil)))

	//SSN will become active on next cycle
	p.Zproxy.UpdateWallet(sdk.Cfg.VerifierKey)

	//we need to increase blocknum, in order to Gzil won't mint anything. Really minting is over.
	sdk.IncreaseBlocknum(10)
	check(p.Zproxy.AssignStakeReward(sdk.Cfg.AzilSsnAddress, sdk.Cfg.AzilSsnRewardShare))
}

func check(tx *transaction.Transaction, err error) (*transaction.Transaction, error) {
	if err != nil {
		_, file, no, _ := runtime.Caller(1)
		log.Fatal("TRANSACTION FAILED, " + file + ":" + strconv.Itoa(no))
	}
	return tx, err
}

const qa = "000000000000"

func zil(amount int) string {
	if amount == 0 {
		return "0"
	}
	return fmt.Sprintf("%d%s", amount, qa)
}