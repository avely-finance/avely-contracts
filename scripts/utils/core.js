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

const { BN, bytes, units } = require('@zilliqa-js/util');
const { Zilliqa } = require('@zilliqa-js/zilliqa');
const { toBech32Address, getAddressFromPrivateKey } = require('@zilliqa-js/crypto');
const expect = require('@truffle/expect');

//parse input arguments
//https://github.com/yargs/yargs/blob/main/docs/api.md
const argv = require('yargs/yargs')(process.argv.slice(2))
  .default('env', 'local')
  .default('network', 'local')
  .argv;

//read environment constants from file
try {
  require('dotenv-safe').config({path: '.env.' + argv.env});
} catch (err) {
  console.error(err.message);
  process.exit(1);
}

const config = require('../../scripts-config.js');
config.network = argv.network;

class AvelyCore
{
  constructor (options)
  {
    options = options || {};
    expect.options(options, ['networks', 'network']);
    this.logger = options.logger || { log: function () {} };
    if (undefined === options.networks[options.network]) {
      throw new Error('Unknown network: ' + options.network);
    }
    this.network = options.networks[options.network];
    this.zilliqa = new Zilliqa(this.network.endPoint);
    this.msgVersion = 1; // current msgVersion
    this.VERSION = bytes.pack(this.network.chainId, this.msgVersion);

    // Populate the wallet with an account
    this.zilliqa.wallet.addByPrivateKey(this.network.fromPK);
    this.fromAddress = getAddressFromPrivateKey(this.network.fromPK);
    console.log(`My account address is: ${this.fromAddress}`);
    console.log(`My account bech32 address is: ${toBech32Address(this.fromAddress)}`);
  }

  async getBalance()
  {
    const response = await this.zilliqa.blockchain.getBalance(this.fromAddress);
    return response.result.balance;
  }

  async isGasSufficient()
  {
    // Get Minimum Gas Price from blockchain
    const minGasPrice = await this.zilliqa.blockchain.getMinimumGasPrice();
    console.log(`Current Minimum Gas Price: ${minGasPrice.result}`);
    this.gasPriceQa = units.toQa(this.network.gasPrice, units.Units.Li); // Gas Price that will be used by all transactions
    console.log(`My Gas Price ${this.gasPriceQa.toString()}`);
    const isGasSufficient = this.gasPriceQa.gte(new BN(minGasPrice.result)); // Checks if your gas price is less than the minimum gas price

    return isGasSufficient;

  }
}

module.exports = new AvelyCore(config);
