const Deployer = require('../utils/deployer.js');
const Contracts = require('../utils/contracts.js');

(async() => {
  try {
    const Buffer = await Contracts.newFromFile('buffer', []);

    await Deployer.deploy(Buffer);

  } catch (err) {
    console.error(err.message);
  }
})();
