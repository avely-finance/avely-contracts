const Deployer = require('../utils/deployer.js');
const Contracts = require('../utils/contracts.js');

(async() => {
  try {
    const AZil = await Contracts.newFromFile('aZil');

    await Deployer.deploy(AZil);

  } catch (err) {
    console.error(err.message);
  }
})();
