# Multisig Wallet

## Commands

1. Submit set holder transaction

```sh
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitSetHolderAddressTransaction --addr=0x166862bdd5d76b3a4775d2494820179d582acac5
```

2. Sync buffers

```sh
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitChangeBuffersTransaction
```

3. Upnause In, Out and ZRC2

```sh
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitUnPauseInTransaction
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitUnPauseOutTransaction
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitUnPauseZrc2Transaction
```
