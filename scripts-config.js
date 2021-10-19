module.exports = {
  networks: {
    local: {
      endPoint: 'http://zilliqa_server:5555',
      chainId: 222,
      gasPrice: 2000, //Gas Price that will be used by all transactions
      gasLimitContractDeploy: 10000, //Gas limit for contract deploy
      gasLimitContractCall: 8000, //Gas limit for contract call
      gasLimitPay: 50, //Gas limit for payment transaction
      fromPK: process.env.ZIL_PRIVATE_KEY1 // Account PK to send txs from
    },
    testnet: {
      endPoint: '...',
      chainId: 333,
      gasPrice: 2000, //Gas Price that will be used by all transactions
      gasLimitContractDeploy: 1, //Gas limit for contract deploy
      gasLimitContractCall: 1, //Gas limit for contract call
      gasLimitPay: 1, //Gas limit for payment transaction
      fromPK: '' // Account PK to send txs from
    }
  }
};
