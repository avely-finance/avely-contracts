//  Copyright (C) 2018 Zilliqa
//
//  This file is part of Zilliqa-Javascript-Library.
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with this program.  If not, see <https://www.gnu.org/licenses/>.


require('dotenv-safe').config();

process.on('uncaughtException', err => {
  console.error(err.message)
  process.exit(1)
})

const { argv, env } = require('process');
if (argv.length < 3) {
  throw new Error("Format is: yarn run deploy contractName");
}
const contractCode = require('./utils.js').readContractCodeByName(argv[2]);

const { BN, Long, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');
const zilliqa = new Zilliqa(env.ZIL_HOST);

// These are set by the core protocol, and may vary per-chain.
// You can manually pack the bytes according to chain id and msg version.
// For more information: https://apidocs.zilliqa.com/?shell#getnetworkid

const chainId = env.ZIL_CHAIN_ID;
const msgVersion = 1; // current msgVersion
const VERSION = bytes.pack(chainId, msgVersion);

// Populate the wallet with an account
zilliqa.wallet.addByPrivateKey(env.ZIL_PRIVATE_KEY1);
const address = getAddressFromPrivateKey(env.ZIL_PRIVATE_KEY1);
console.log(`My account address is: ${address}`);
console.log(`My account bech32 address is: ${toBech32Address(address)}`);

async function deployContract() {
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

    // Deploy a contract
    console.log(`Deploying a new contract....`);

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
    const contract = zilliqa.contracts.new(contractCode, init);
    // Deploy the contract.
    // Also notice here we have a default function parameter named toDs as mentioned above.
    // A contract can be deployed at either the shard or at the DS. Always set this value to false.
    const [deployTx, deployedContract] = await contract.deployWithoutConfirm(
      {
        version: VERSION,
        gasPrice: txGasPrice,
        gasLimit: Long.fromNumber(env.TX_DEPLOY_GAS_LIMIT),
      },
      false,
    );

    // process confirm
    console.log(`The transaction id is:`, deployTx.id);
    console.log(`Waiting transaction be confirmed`);
    const confirmedTxn = await deployTx.confirm(deployTx.id);

    console.log(`The transaction status is:`);
    console.log(confirmedTxn.receipt);
    if (confirmedTxn.receipt.success === true) {
      console.log(`Contract address is: ${deployedContract.address}`);
    }
    
    //Get the contract state
    console.log('Getting contract state...');
    const state = await deployedContract.getState();
    console.log('The state of the contract is:');
    console.log(JSON.stringify(state, null, 4));
  } catch (err) {
    console.log(err);
  }
}

deployContract();
