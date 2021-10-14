require('dotenv-safe').config();
const { env } = require('process');
const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const {
  toBech32Address,
  getAddressFromPrivateKey,
} = require('@zilliqa-js/crypto');

const zilliqa = new Zilliqa(env.ZIL_HOST);

// These are set by the core protocol, and may vary per-chain.
// You can manually pack the bytes according to chain id and msg version.
// For more information: https://apidocs.zilliqa.com/?shell#getnetworkid

const chainId = env.ZIL_CHAIN_ID;
const msgVersion = 1; // current msgVersion
const VERSION = bytes.pack(chainId, msgVersion);

// Populate the wallet with an account
zilliqa.wallet.addByPrivateKey(env.ZIL_PRIVATE_KEY1);

const testRcptAddress = '';

const address = getAddressFromPrivateKey(env.ZIL_PRIVATE_KEY1);
console.log(`My account address is: ${address}`);
console.log(`My account bech32 address is: ${toBech32Address(address)}`);

async function testBlockchain() {
  try {
    // Get Balance
    const balance = await zilliqa.blockchain.getBalance(address);
    // Get Minimum Gas Price from blockchain
    const minGasPrice = await zilliqa.blockchain.getMinimumGasPrice();

    // Account balance (See note 1)
    console.log(`Your account balance is:`);
    console.log(balance.result);
    console.log(`Current Minimum Gas Price: ${minGasPrice.result}`);
    const txGasPrice = units.toQa(env.TX_GAS_PRICE, units.Units.Li); // Gas Price that will be used by all transactions
    console.log(`My Gas Price ${txGasPrice.toString()}`);
    const isGasSufficient = txGasPrice.gte(new BN(minGasPrice.result)); // Checks if your gas price is less than the minimum gas price
    console.log(`Is the gas price sufficient? ${isGasSufficient}`);

    // Send a transaction to the network
    console.log('Sending a payment transaction to the network...');
    const tx = await zilliqa.blockchain.createTransaction(
      // Notice here we have a default function parameter named toDs which means the priority of the transaction.
      // If the value of toDs is false, then the transaction will be sent to a normal shard, otherwise, the transaction.
      // will be sent to ds shard. More info on design of sharding for smart contract can be found in.
      // https://blog.zilliqa.com/provisioning-sharding-for-smart-contracts-a-design-for-zilliqa-cd8d012ee735.
      // For payment transaction, it should always be false.
      zilliqa.transactions.new(
        {
          version: VERSION,
          toAddr: env.ZIL_TEST_RECIPIENT_ADDRESS1,
          amount: new BN(units.toQa('1', units.Units.Zil)), // Sending an amount in Zil (1) and converting the amount to Qa
          gasPrice: txGasPrice, // Minimum gasPrice veries. Check the `GetMinimumGasPrice` on the blockchain
          gasLimit: Long.fromNumber(env.TX_PAY_GAS_LIMIT),
        },
        false,
      ),
    );

    console.log(`The transaction status is:`);
    console.log(tx.receipt);


    // Deploy a contract
    console.log(`Deploying a new contract....`);
    const code = `scilla_version 0

    (* HelloWorld contract *)

    import ListUtils

    (***************************************************)
    (*               Associated library                *)
    (***************************************************)
    library HelloWorld

    let not_owner_code = Int32 1
    let set_hello_code = Int32 2

    (***************************************************)
    (*             The contract definition             *)
    (***************************************************)

    contract HelloWorld
    (owner: ByStr20)

    field welcome_msg : String = ""

    transition setHello (msg : String)
      is_owner = builtin eq owner _sender;
      match is_owner with
      | False =>
        e = {_eventname : "setHello()"; code : not_owner_code};
        event e
      | True =>
        welcome_msg := msg;
        e = {_eventname : "setHello()"; code : set_hello_code};
        event e
      end
    end


    transition getHello ()
        r <- welcome_msg;
        e = {_eventname: "getHello()"; msg: r};
        event e
    end`;

    const init = [
      // this parameter is mandatory for all init arrays
      {
        vname: '_scilla_version',
        type: 'Uint32',
        value: '0',
      },
      {
        vname: 'owner',
        type: 'ByStr20',
        value: `${address}`,
      },
    ];

    // Instance of class Contract
    const contract = zilliqa.contracts.new(code, init);

    // Deploy the contract.
    // Also notice here we have a default function parameter named toDs as mentioned above.
    // A contract can be deployed at either the shard or at the DS. Always set this value to false.
    const [deployTx, hello] = await contract.deploy(
      {
        version: VERSION,
        gasPrice: txGasPrice,
        gasLimit: Long.fromNumber(env.TX_DEPLOY_GAS_LIMIT),
      },
      33,
      1000,
      false,
    );

    // Introspect the state of the underlying transaction
    console.log(`Deployment Transaction ID: ${deployTx.id}`);
    console.log(`Deployment Transaction Receipt:`);
    console.log(deployTx.txParams.receipt);

    // Get the deployed contract address
    console.log('The contract address is:');
    console.log(hello.address);

    //Following line added to fix issue https://github.com/Zilliqa/Zilliqa-JavaScript-Library/issues/168
    const deployedContract = zilliqa.contracts.at(hello.address);

    // Create a new timebased message and call setHello
    // Also notice here we have a default function parameter named toDs as mentioned above.
    // For calling a smart contract, any transaction can be processed in the DS but not every transaction can be processed in the shards.
    // For those transactions are involved in chain call, the value of toDs should always be true.
    // If a transaction of contract invocation is sent to a shard and if the shard is not allowed to process it, then the transaction will be dropped.
    const newMsg = 'Hello, the time is ' + Date.now();
    console.log('Calling setHello transition with msg: ' + newMsg);
    const callTx = await hello.call(
      'setHello',
      [
        {
          vname: 'msg',
          type: 'String',
          value: newMsg,
        },
      ],
      {
        // amount, gasPrice and gasLimit must be explicitly provided
        version: VERSION,
        amount: new BN(0),
        gasPrice: txGasPrice,
        gasLimit: Long.fromNumber(env.TX_CALL_GAS_LIMIT),
      },
      33,
      1000,
      false,
    );

    // Retrieving the transaction receipt (See note 2)
    console.log(JSON.stringify(callTx.receipt, null, 4));

    //Get the contract state
    console.log('Getting contract state...');
    const state = await deployedContract.getState();
    console.log('The state of the contract is:');
    console.log(JSON.stringify(state, null, 4));
  } catch (err) {
    console.log(err);
  }
}

testBlockchain();
