# Tools

- [Configuration](#configuration)
- [Useful Links](#useful-links)
- [User commands](#user-commands)
- [Admin commands](#admin-commands)
    - [Deploy](#deploy)
    - [Utils, information](#utils-information)
    - [New reward cycle](#new-reward-cycle)
    - [Withdrawals](#withdrawals)
    - [Swap requests \(part of chown stake process\)](#swap-requests-part-of-chown-stake-process)

## Configuration

Create and configure `.env.testnet` and `.env.mainnet`. Ask the team for private keys for shared admin and user accounts

## Useful Links

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

## User commands

1. go run tools/user_cmd.go --chain=testnet --cmd=delegate --usr=1 --amount=100

## Admin commands

### Deploy

1. Deploy basic contracts

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=deploy
```

2. Copy-paste new contract addresses to the `config.json`

3. Deploy second buffer

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=deploy_buffer
```

4. Update `config.json` again and run

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=sync_buffers
```

5. Unpause in/out/all AZil

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=unpause_in
```

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=unpause_out
```

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=unpause_all
```


6. Init Holder with min stake

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=init_holder
```

### Utils, information

1. Convert address from bech32 to base16

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=from_bech32 --addr=<bech32 addr>
```

2. Convert address from base16 to bech32

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=to_bech32 --addr=<base16 addr>
```

3. Show transaction

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=show_tx --addr=<transaction hash>
```

4. Get Active Buffer

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=get_active_buffer
```

5. Show Stake Rewards on the main staking contract

```sh
$ go run tools/admin_cmd.go --chain=mainnet --cmd=show_rewards --ssn=0x2afe9e18EdD39D927d0FffF8990612FC4aFa2295 --addr=0x30B5259a4E89Dc12B6da7883A9D3cd691F03b386
```

### New reward cycle

1. Drain Buffer

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=drain_buffer --addr=<buffer addr>
```

2. ReDelegate stakes after swap requests confirmation (show-only mode)

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=show_redelegate
```

3. ReDelegate stakes after swap requests confirmation

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=redelegate
```

4. Perform Autorestake

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=autorestake
```

### Withdrawals

1. Show blocks with withdrawals, ready for claim

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=show_claim_withdrawal
```

2. Claim withdrawals

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=claim_withdrawal
```

### Swap requests (part of chown stake process)

1. Show swap request(s)

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=show_swap_requests
```

2. Confirm swap request(s)

```sh
$ go run tools/admin_cmd.go --chain=testnet --cmd=confirm_swap_requests
```
