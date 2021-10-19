function readContractCodeByName (name, ext = '.scilla')
{
  const { sep, resolve } = require('path');
  const { existsSync, readFileSync } = require('fs');
  const contractsPath = resolve(__dirname + sep + '..' + sep + '..' + sep + 'contracts' + sep);
  const fullContractPath = contractsPath + sep + (name.endsWith(ext) ? name : name + ext);
  if (!existsSync(fullContractPath)) {
    throw new Error('Contract not found at path: ' + fullContractPath);
  }
  return readFileSync(fullContractPath).toString();
}

module.exports = { readContractCodeByName };
