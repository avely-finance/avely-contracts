const Deployer = require('../utils/deployer.js');
const Contracts = require('../utils/contracts.js');

(async() => {
  try {
    const Holder = await Contracts.newFromFile('holder', []);

    await Deployer.deploy(Holder);

  } catch (err) {
    console.error(err.message);
  }
})();
