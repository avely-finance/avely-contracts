# Testnet contracts

## Preabmle

In past for testing contracts in testnet we were using original Zilliqa's contracts, previously deployed by Zilliqa.
In that setup reward cycle change transition supposed to be called by [Verifier](https://dev.zilliqa.com/docs/staking/phase1/staking-phase1-overview/) (part of the system, responsible for data verification).

Disadvantages of such aprroach:

* We weren't control zilliqa's contracts
* We should wait some time until Verifier changes reward cycle (about [75 Final Blocks](https://dev.zilliqa.com/docs/staking/phase1/staking-general-information/#testnet), ~30 mins).

It was annoying, for example, during withdrawals testing.
Now we [redeployed](https://github.com/avely-finance/avely-contracts/blob/main/docs/tools.md#deploy-zilliqa-staking-contracts) copies of all zilliqa'a staking contracts ([gzil](https://github.com/Zilliqa/staking-contract/blob/main/contracts/gzil.scilla),
[proxy](https://github.com/Zilliqa/staking-contract/blob/main/contracts/proxy.scilla),
[ssnlist](https://github.com/Zilliqa/staking-contract/blob/main/contracts/ssnlist.scilla)), and using them in our test setup.
This way we can change reward cycle ourselves.

## How to change reward cycle

* clone project `git clone git@github.com:avely-finance/avely-contracts.git`
* create .env.testnet (change [.env.example](https://github.com/avely-finance/avely-contracts/blob/main/.env.example)) and put sensitive information there (ask team)
* execute
```sh
$ go run tools/zilliqa_staking_cmd.go --chain=testnet --cmd=next_cycle
```

