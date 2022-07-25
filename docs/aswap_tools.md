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
     - [Claim Owner](#claim-owner)
- [User commands](#user-commands)
     - [Add Liquidity](#add-liquidity)
     - [Remove Liquidity](#remove-liquidity)
     - [Swap exact ZIL for tokens](#swap-exact-zil-for-tokens)
     - [Swap exact tokens for ZIL](#swap-exact-tokens-for-zil)

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

Command will call `SetLiquidityFee(new_fee: Uint256)` transition

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=set_liquidity_fee --value=9000
```

### Set Treasury Fee

Command will call `SetTreasuryFee(new_fee: Uint128)` transition

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=set_treasury_fee --value=550
```

### Set Treasury Address

Command will call `SetTreasuryAddress(new_address: ByStr20)` transition

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=set_treasury_address --value=0x0000000000000000000000000000000000000000
```

### Change Owner

Command will call `ChangeOwner(new_owner: ByStr20)` transition

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=change_owner --value=0x0000000000000000000000000000000000000000
```

### Claim Owner

Command will call `ClaimOwner()` transition

```sh
$ go run tools/aswap_admin_cmd.go --chain=testnet --cmd=claim_owner
```

## User commands

### Add Liquidity

Command will call `AddLiquidity()` transition

```sh
$ go run tools/aswap_user_cmd.go help add_liquidity
NAME:
   aswap-user-cmd add_liquidity - add liquidity

USAGE:
   aswap-user-cmd add_liquidity [command options] [arguments...]

OPTIONS:
   --chain value                    Chain (default: "local")
   --deadline_block value           Deadline block for localnet, gap for other (default: 0)
   --max_token_amount value         Maximum token amount (default: 0)
   --min_contribution_amount value  Minimum contribution amount (default: 0)
   --token_address value            Token address, StZil by default
   --zil_amount value               Native ZIL _amount to send (default: 0)
```

Example

```sh
$ go run tools/aswap_user_cmd.go add_liquidity --chain=testnet --zil_amount=100 --min_contribution_amount=10 --max_token_amount=100 --deadline_block=56789
```

### Remove Liquidity

Command will call `RemoveLiquidity()` transition

```sh
$ go run tools/aswap_user_cmd.go help remove_liquidity
NAME:
   aswap-user-cmd remove_liquidity - remove liquidity

USAGE:
   aswap-user-cmd remove_liquidity [command options] [arguments...]

OPTIONS:
   --chain value                Chain (default: "local")
   --contribution_amount value  Contribution amount (default: 0)
   --deadline_block value       Deadline block for localnet, gap for other (default: 0)
   --min_token_amount value     Minimum token amount (default: 0)
   --min_zil_amount value       Minimum ZIL amount (default: 0)
   --token_address value        Token address, StZil by default
```

Example

```sh
$ go run tools/aswap_user_cmd.go remove_liquidity --chain=testnet --contribution_amount=100 --min_zil_amount=100 --min_token_amount=100 --deadline_block=56789
```

### Swap exact ZIL for tokens

Command will call `SwapExactZILForTokens` transition

```sh
$ go run tools/aswap_user_cmd.go help swap_exact_zil_for_tokens
NAME:
   aswap-user-cmd swap_exact_zil_for_tokens - Swaps exact zil for tokens

USAGE:
   aswap-user-cmd swap_exact_zil_for_tokens [command options] [arguments...]

OPTIONS:
   --chain value              Chain (default: "local")
   --deadline_block value     Deadline block for localnet, gap for other (default: 0)
   --min_token_amount value   Minimum token amount (default: 0)
   --recipient_address value  Recipient address
   --token_address value      Token address, StZil by default
   --zil_amount value         Native ZIL _amount to send (default: 0)
```

Example

```sh
$ go run tools/aswap_user_cmd.go swap_exact_zil_for_tokens --chain=testnet --zil_amount=100 --min_token_amount=100 --deadline_block=56789 --recipient_address=0x0000000000000000000000000000000000000001
```

### Swap exact tokens for ZIL

Command will call `SwapExactTokensForZIL` transition

```sh
$ go run tools/aswap_user_cmd.go help swap_exact_tokens_for_zil
NAME:
   aswap-user-cmd swap_exact_tokens_for_zil - Swaps exact tokens for zil

USAGE:
   aswap-user-cmd swap_exact_tokens_for_zil [command options] [arguments...]

OPTIONS:
   --chain value              Chain (default: "local")
   --deadline_block value     Deadline block for localnet, gap for other (default: 0)
   --min_zil_amount value     Minimum ZIL amount (default: 0)
   --recipient_address value  Recipient address
   --token_address value      Token address, StZil by default
   --token_amount value       Token amount (default: 0)

```

Example

```sh
$ go run tools/aswap_user_cmd.go swap_exact_tokens_for_zil --chain=testnet --token_amount=100 --min_zil_amount=100 --deadline_block=56789 --recipient_address=0x0000000000000000000000000000000000000001
```
