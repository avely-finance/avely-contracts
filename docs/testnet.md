# Testnet contracts

<!-- MarkdownTOC -->

- [Preabmle](#preabmle)
- [How to change reward cycle](#how-to-change-reward-cycle)
- [Current setup](#current-setup)
    - [stZIL contract/token](#stzil-contracttoken)
    - [SSNs](#ssns)
    - [Multisig contract](#multisig-contract)
        - [Multisig owner 1](#multisig-owner-1)
        - [Multisig owner 2](#multisig-owner-2)

<!-- /MarkdownTOC -->


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

## Current setup

### stZIL contract/token

[zil1hkq4uugau2qea9ehtr8r8z2xdvrez5xpfz8tl0](https://viewblock.io/zilliqa/address/zil1hkq4uugau2qea9ehtr8r8z2xdvrez5xpfz8tl0?network=testnet) / 0xbd815e711de2819e973758ce3389466b079150c1

### SSNs

* Nodamatics zil1w0ckw5l2jlrhqgql66hqgpenfrzgxgul05g4te / 0x73f16753ea97C770201Fd6Ae04073348c483239f
* Zillacracy zil14xum8cvuhsfzh9upgfd39cracwt2zkzsak4ymr / 0xA9B9B3e19cbC122B9781425B12e07dC396A15850

### Multisig contract

[zil1l2rtpzcvwdthts2ltvutyau9ptwpskef7pvgkd](https://viewblock.io/zilliqa/address/zil1l2rtpzcvwdthts2ltvutyau9ptwpskef7pvgkd?network=testnet) / 0xfa86b08b0c735775c15f5b38b277850adc185b29

#### Multisig owner 1

zil17c2804u3j3uwttl7r77e5rxuam5lmepdf2l87e / 0xf61477D7919478e5AfFe1fbd9A0CDCeee9fdE42d

#### Multisig owner 2

zil1unaeqx3050y8ae5pe6ah6ujk24lspvq4rnec77 / 0xE4fB901A2FA3C87ee681cEbb7D7256557f00b015
