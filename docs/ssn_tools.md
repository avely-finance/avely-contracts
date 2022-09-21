# SSN contract tools

<!-- MarkdownTOC -->

- [Configuration](#configuration)
- [Roles](#roles)
- [Admin commands](#admin-commands)
     - [Deploy](#deploy)
     - [Print state](#print-state)
     - [Change Owner](#change-owner)
     - [Claim Owner](#claim-owner)
     - [Change Zproxy Addres](#change-zproxy-addres)
     - [Update Receiving Address](#update-receiving-address)
     - [Update Comission](#update-comission)
     - [Withdraw Comission](#withdraw-comission)

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
$ go run tools/ssn_cmd.go --chain=testnet --cmd=deploy
```
After successful deploy put deployed contract address to config.js and/or .env.your-chain

See `StZilSsnAddress` keyname.

### Print state

```sh
$ go run tools/ssn_cmd.go --chain=testnet --cmd=print_state
...
INFO[2022-09-17T16:53:31Z] {
     "_balance": "0",
     "owner": "0xf61477d7919478e5affe1fbd9a0cdceee9fde42d",
     "staging_owner": {
          "argtypes": [
               "ByStr20"
          ],
          "arguments": [],
          "constructor": "None"
     },
     "zproxy": "0x310dbc947ae5250b7f84247eb52a4fdb85b5e35a"
}
```

### Change Owner

Command will call `ChangeOwner(new_owner: ByStr20)` transition

```sh
$ go run tools/ssn_cmd.go --chain=testnet --cmd=change_owner --recipient=0x0000000000000000000000000000000000000000
```

### Claim Owner

Command will call `ClaimOwner()` transition

```sh
$ go run tools/ssn_cmd.go --chain=testnet --cmd=claim_owner --usr=1
```

Possible values for --usr parametes is 1, 2, 3, which will set actor to config.Key1, config.Key2, config.Key3 correspondingly;
--usr=admin will set actor to AdminKey, --usr=owner will set actor to OwnerKey

### Change Zproxy Addres

Command will call `ChangeZproxy(new_addr: ByStr20)` transition

```sh
$ go run tools/ssn_cmd.go --chain=testnet --cmd=change_zproxy --param1=0x0000000000000000000000000000000000000001
```

### Update Receiving Address

Command will call `UpdateReceivingAddr(new_addr: ByStr20)` transition

```sh
$ go run tools/ssn_cmd.go --chain=testnet --cmd=update_receiving_addr --param1=0x0000000000000000000000000000000000000001
```

### Update Comission

Command will call `UpdateComm(new_rate: Uint128)` transition

```sh
$ go run tools/ssn_cmd.go --chain=testnet --cmd=update_comm --param1=12345
```
If you got Int32 -10 SSNNotExists error, that means that SSN was not added to ssnlist contract.

You can [add it](tools.md#deploy) in testnet, or ask ssnlist contract administrators in mainnet.


### Withdraw Comission

Command will call `WithdrawComm()` transition

```sh
$ go run tools/ssn_cmd.go --chain=testnet --cmd=withdraw_comm
```

Error SSNNoComm => Int32 -14 is possible. That means that SSN has no commission yet.
