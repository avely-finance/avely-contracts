# ASwap contract tools

<!-- MarkdownTOC -->

- [Configuration](#configuration)
- [Roles](#roles)
- [Admin commands](#admin-commands)
    - [Deploy](#deploy)
    - [Print state](#print-state)
    - [Set Liquidity Fee](#set-liquidity-fee)
    - [Set Treasury Fee](#set-treasury-fee)
    - [Set Treasury Address](#set-treasury-address)
    - [Change Owner](#change-owner)

<!-- /MarkdownTOC -->


## Configuration

Create and configure `.env.testnet` and `.env.mainnet`. Ask the team for private keys for shared admin and user accounts.

## Roles

1. *Admin*. See ADMIN/ADMIN_KEY in [.env.example](../.env.example) and Admin/AdminKey in [config.json](../config.json)
1. *Owner*. See OWNER/OWNER_KEY in [.env.example](../.env.example) and Owner/OwnerKey in [config.json](../config.json)

## Admin commands

### Deploy

Deploy will be done by *Admin*.
Contract owner field (see [field owner: ...](../contracts/aswap.scilla)) will be set to *Owner*.
After deploy all admin transitions will be executed by *Owner*.

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=deploy
```

After successful deploy put deployed contract address to config.js and/or .env.your-chain

### Print state

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=print_state
...
INFO[] {
     "_balance": "0",
     "balances": {},
     "liquidity_fee": "10000",
     "min_lp": "100000000000000",
     "owner": "0x...",
     "pause": "0",
     "pools": {},
     "staging_owner": {
          "argtypes": [
               "ByStr20"
          ],
          "arguments": [],
          "constructor": "None"
     },
     "total_contributions": {},
     "treasury_address": "0x...",
     "treasury_fee": "500"
}

```

### Set Liquidity Fee

Command will call transition `SetLiquidityFee(new_fee: Uint256)`

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=set_liquidity_fee --value=9000
```

### Set Treasury Fee

Command will call transition `SetTreasuryFee(new_fee: Uint128)`

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=set_treasury_fee --value=550
```

### Set Treasury Address

Command will call transition `SetTreasuryAddress(new_address: ByStr20)`

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=set_treasury_address --value=0x0000000000000000000000000000000000000000
```

### Change Owner

Command will call transition `ChangeOwner(new_owner: ByStr20)`

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=change_owner --value=0x0000000000000000000000000000000000000000
```
