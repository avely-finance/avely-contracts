# Testnet

## Configuration

Create and configure `.env.testnet`. Ask the team for private keys for shared admin and user accounts

## Links

Staking testnet Dashboard:
https://testnet-stake.zilliqa.com/dashboard


Admin account:
https://viewblock.io/zilliqa/address/zil17c2804u3j3uwttl7r77e5rxuam5lmepdf2l87e?network=testnet


## Admin commands

1. go run tools/admin_actions.go --chain=testnet --cmd=deploy
2. go run tools/admin_actions.go --chain=testnet --cmd=deploy_buffer
3. go run tools/admin_actions.go --chain=testnet --cmd=sync_buffers


## User commands

1. go run tools/user_actions.go --chain=testnet --cmd=delegate --usr=1 --amount=10
