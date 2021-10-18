const Deployer = require('./utils/deployer.js');
const Contracts = require('./utils/contracts.js');

(async() => {
  try {

    //contract 1 with default parameters
    const MyContract1 = await Contracts.newFromFile('helloWorld');

    //contract 2 with extended parameters
    const extendedInit = [{
      vname: 'foo',
      type: 'Uint32',
      value: '555'
    }];
    const MyContract2 = await Contracts.newFromFile('fooBar', extendedInit);

    //deploy
      await Deployer.deploy(MyContract1);
      await Deployer.deploy(MyContract2);

  } catch (err) {
    console.error(err.message)
  }
})();


/*
if (argv.length < 3) {
  throw new Error("Format is: yarn run deploy contractName");
}*/
