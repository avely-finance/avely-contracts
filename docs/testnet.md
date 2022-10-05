# Testnet contracts

<!-- MarkdownTOC -->

- [Preabmle](#preabmle)
- [How to change reward cycle](#how-to-change-reward-cycle)
- [Roles](#roles)
- [Current protocol setup](#current-protocol-setup)
    - [Whitelisted SSNs](#whitelisted-ssns)
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
$ go run tools/zilliqa_staking_cmd.go --chain=testnet --cmd=next_cycle --amount=100
```

## Roles

Admin, owner, verifier: 0xf61477D7919478e5AfFe1fbd9A0CDCeee9fdE42d

## Current protocol setup

* ASwap: 0x7de8d9ce9fa4725b7e6811d96fb690d7d7f86b42 (owner is *Admin*) zil10h5dnn5l53e9klngz8vkld5s6ltls66zs4kzn7
* ASwap: 0xd93b636ecd4e6be8df62c0c0197b728fa4095111 (owner is *Multisig*) zil1myakxmkdfe473hmzcrqpj7mj37jqj5g3uguq5s
* Proxy: 0x310dbc947ae5250b7f84247eb52a4fdb85b5e35a
* SsnList: 0x00cab6f10a801622c481ce7bf0737095fba39417
* StZil: 0xbd815e711de2819e973758ce3389466b079150c1 [zil1hkq4uugau2qea9ehtr8r8z2xdvrez5xpfz8tl0](https://viewblock.io/zilliqa/address/zil1hkq4uugau2qea9ehtr8r8z2xdvrez5xpfz8tl0?network=testnet)
* Buffers:
  * 0x73bb9c1abb055a43e058d4711ad1fe6bccb64269,
  * 0x4439ef882cc395e4e55f557aa5b8908d14aae764,
  * 0x27be8d0a99774c32ff10e1f90d63f2705a111850
* Holder: 0x6c242e56a377f69bd44db49aed2124c8393b5054,

### Whitelisted SSNs

* Nodamatics zil1w0ckw5l2jlrhqgql66hqgpenfrzgxgul05g4te / 0x73f16753ea97C770201Fd6Ae04073348c483239f
* Zillacracy zil14xum8cvuhsfzh9upgfd39cracwt2zkzsak4ymr / 0xA9B9B3e19cbC122B9781425B12e07dC396A15850

### Multisig contract

* [zil1l2rtpzcvwdthts2ltvutyau9ptwpskef7pvgkd](https://viewblock.io/zilliqa/address/zil1l2rtpzcvwdthts2ltvutyau9ptwpskef7pvgkd?network=testnet)
* 0xfa86b08b0c735775c15f5b38b277850adc185b29
* [**Multisig front-end**](https://avely-multisig.web.app/#/login)

#### Multisig owner 1

zil17c2804u3j3uwttl7r77e5rxuam5lmepdf2l87e / 0xf61477D7919478e5AfFe1fbd9A0CDCeee9fdE42d

#### Multisig owner 2

zil1unaeqx3050y8ae5pe6ah6ujk24lspvq4rnec77 / 0xE4fB901A2FA3C87ee681cEbb7D7256557f00b015
