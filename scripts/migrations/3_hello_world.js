const Deployer = require('../utils/deployer.js');
const Contracts = require('../utils/contracts.js');

(async() => {
  try {
    const hellowWorld = await Contracts.newFromFile('helloWorld');

    await Deployer.deploy(hellowWorld);

  } catch (err) {
    console.error(err.message);
  }
})();
