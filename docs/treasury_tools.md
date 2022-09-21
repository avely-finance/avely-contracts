# Treasury contract tools

<!-- MarkdownTOC -->

- [Configuration](#configuration)
- [Roles](#roles)
- [Admin commands](#admin-commands)
     - [Deploy](#deploy)
     - [Print state](#print-state)
     - [Withdraw](#withdraw)
     - [Change Owner](#change-owner)
     - [Claim Owner](#claim-owner)

<!-- /MarkdownTOC -->

## Configuration

Create and configure `.env.testnet` and `.env.mainnet`. Ask the team for private keys for shared admin and user accounts.

## Roles

1. *Admin*. See ADMIN/ADMIN_KEY in [.env.example](../.env.example) and Admin/AdminKey in [config.json](../config.json)
1. *Owner*. See OWNER/OWNER_KEY in [.env.example](../.env.example) and Owner/OwnerKey in [config.json](../config.json)

## Admin commands

### Deploy

Deploy will be done by *Admin*.
Contract owner field (see [field owner: ...](../contracts/treasury.scilla)) will be set to *Owner*.
After deploy all admin transitions will be executed by *Owner*.

```sh
$ go run tools/treasury_admin_cmd.go --chain=testnet --cmd=deploy
```

After successful deploy put deployed contract address to config.js and/or .env.your-chain

### Print state

```sh
$ go run tools/treasury_admin_cmd.go --chain=testnet --cmd=print_state
...
INFO[2022-09-17T16:53:31Z] {
     "_balance": "980000000001000",
     "owner": "0x6cd3667ba79310837e33f0aecbe13688a6cbca32",
     "staging_owner": {
          "argtypes": [
               "ByStr20"
          ],
          "arguments": [],
          "constructor": "None"
     }
}
```

### Withdraw

Command will call `Withdraw(recipient: ByStr20, amount: Uint128)` transition

```sh
$ go run tools/treasury_admin_cmd.go --chain=testnet --cmd=withdraw --recipient=0x.... --value=123
```

### Change Owner

Command will call `ChangeOwner(new_owner: ByStr20)` transition

```sh
$ go run tools/treasury_admin_cmd.go --chain=testnet --cmd=change_owner --recipient=0x0000000000000000000000000000000000000000
```

### Claim Owner

Command will call `ClaimOwner()` transition

```sh
$ go run tools/treasury_admin_cmd.go --chain=testnet --cmd=claim_owner --usr=1
```

Possible values for --usr parametes is 1, 2, 3, which will set actor to config.Key1, config.Key2, config.Key3 correspondingly;
--usr=admin will set actor to AdminKey, --usr=owner will set actor to OwnerKey
