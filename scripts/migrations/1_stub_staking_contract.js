const Deployer = require('../utils/deployer.js');
const Contracts = require('../utils/contracts.js');

(async() => {
  try {
    const StubStakingContract = await Contracts.newFromFile('stubStakingContract');

    await Deployer.deploy(StubStakingContract);

  } catch (err) {
    console.error(err.message);
  }
})();
