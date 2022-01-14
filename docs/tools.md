# Tools

## Configuration

Create and configure `.env.testnet` and `.env.mainnet`. Ask the team for private keys for shared admin and user accounts

## Usefull Links

Staking testnet Dashboard:
https://testnet-stake.zilliqa.com/dashboard

Faucet:
https://dev-wallet.zilliqa.com/faucet?network=testnet

Staking contract addresses:
* https://dev.zilliqa.com/docs/staking/phase1/staking-general-information/
* https://testnet-viewer.zilliqa.com/#staking-contract

Staked Seed Nodes:
https://testnet-viewer.zilliqa.com/

Admin account:
https://viewblock.io/zilliqa/address/zil17c2804u3j3uwttl7r77e5rxuam5lmepdf2l87e?network=testnet

## Admin commands

### Deploy

1. Deploy basic contracts

```sh
$ go run tools/admin_actions.go --chain=testnet --cmd=deploy
```

2. Copy-paste new contract addresses to the `config.json`

3. Deploy second buffer

```sh
$ go run tools/admin_actions.go --chain=testnet --cmd=deploy_buffer
```

4. Update `config.json` again and run

```sh
$ go run tools/admin_actions.go --chain=testnet --cmd=sync_buffers
```

5. Unpause AZil

```sh
$ go run tools/admin_actions.go --chain=testnet --cmd=unpause
```

### Get Info

1. go run tools/admin_actions.go --chain=mainnet --cmd=show_rewards --ssn=0x2afe9e18EdD39D927d0FffF8990612FC4aFa2295 --addr=0x30B5259a4E89Dc12B6da7883A9D3cd691F03b386

## User commands

1. go run tools/user_actions.go --chain=testnet --cmd=delegate --usr=1 --amount=100
