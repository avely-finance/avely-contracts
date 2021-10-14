module.exports = {
  readContractCodeByName: function (name) {
    const { env } = require('process');
    const { sep, resolve } = require('path');
    const { existsSync, readFileSync } = require('fs');

    const contractsPath = resolve((env.PROJECT_CWD || env.INIT_CWD) + sep + 'contracts' + sep);
    const fullContractPath = contractsPath + sep + (name.endsWith('.scilla') ? name : name + '.scilla');
      if (!existsSync(fullContractPath)) {
      throw new Error("Contract not found at path: " + fullContractPath);
    }
    return readFileSync(fullContractPath).toString();
  }
}
