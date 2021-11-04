const Deployer = require('../utils/deployer.js');
const Contracts = require('../utils/contracts.js');

(async() => {
  try {
    const extendedInit = [
      {
        vname: 'azil_ssn_addrress',
        type: 'ByStr20',
        value: '0x381f4008505e940ad7681ec3468a719060caf796'
      },
      {
        vname: 'init_proxy_staking_contract_address',
        type: 'ByStr20',
        value: '0xc6e4fa9abb99f2b3919990ba194d273fd3e21ac9'
      },
      {
        vname: 'init_holder_address',
        type: 'ByStr20',
        value: '0xb2e2c996e6068f4ae11c4cc2c6a189b774819f79'
      }
    ];

    const Buffer = await Contracts.newFromFile('buffer', extendedInit);

    await Deployer.deploy(Buffer);

  } catch (err) {
    console.error(err.message);
  }
})();
