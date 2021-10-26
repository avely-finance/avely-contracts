const Deployer = require('../utils/deployer.js');
const Contracts = require('../utils/contracts.js');

(async() => {
  try {
    const extendedInit = [{
      vname: 'azil_ssnaddr',
      type: 'ByStr20',
      value: '0x166862bdd5d76b3a4775d2494820179d582acac5'
    }];

    const AZil = await Contracts.newFromFile('aZil', extendedInit);

    await Deployer.deploy(AZil);

  } catch (err) {
    console.error(err.message);
  }
})();
