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

const AvelyCore = require('./core.js');
const { BN, Long, bytes, units } = require('@zilliqa-js/util');

class AvelyDeployer
{

  constructor()
  {
  }

  async preDeployCheck(contractInstance)
  {
    //do some checks and return errors (reject)

    const balance = await AvelyCore.getBalance();
    console.log('Your account balance is:' + balance);

    const isGasSufficient = await AvelyCore.isGasSufficient();
    console.log(`Is the gas price sufficient? ${isGasSufficient}`);
  }


  async deploy(contractInstance)
  {
    await this.preDeployCheck();

    // Deploy a contract
    console.log(`Deploying a new contract....`);

    // Deploy the contract.
    // Also notice here we have a default function parameter named toDs as mentioned above.
    // A contract can be deployed at either the shard or at the DS. Always set this value to false.
    const [deployTx, deployedContract] = await contractInstance.deployWithoutConfirm(
      {
        version: AvelyCore.VERSION,
        gasPrice: AvelyCore.gasPriceQa,
        gasLimit: Long.fromNumber(AvelyCore.network.gasLimitContractDeploy),
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
  }
}

module.exports = new AvelyDeployer();
