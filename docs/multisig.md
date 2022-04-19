# Multisig Wallet

## Commands

1. Submit set holder transaction

```sh
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitSetHolderAddressTransaction
```

2. Sync buffers

```sh
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitChangeBuffersTransaction
```

3. Add SSN address

```sh
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitAddSSNTransaction --ssn=ssn_address
```

4. Change Protocol Rewards Fee

```sh
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitChangeRewardsFeeTransaction
```

4. Upnause In, Out and ZRC2

```sh
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitUnPauseInTransaction
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitUnPauseOutTransaction
go run tools/multisig_cmd.go --chain=testnet --tag=SubmitUnPauseZrc2Transaction
```
