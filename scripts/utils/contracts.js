const AvelyCore  = require('./core.js');
const { readContractCodeByName } = require('./fs.js');

class AvelyContracts
{
    constructor()
    {

    }
    
    extendInitData (base, other)
    {
      /*
      input parameters format is:
      [
        {
          vname: '_scilla_version',
          type: 'Uint32',
          value: '0',
        },
        {
          vname: 'owner',
          type: 'ByStr20',
          value: 'aaa',
        },
      ]
      */
      if (!Array.isArray(base) || !Array.isArray(other)) {
          throw new Error('Input parameters must be arrays');
      }
      const result = [];
      base.forEach(function(item, index){
          result[item.vname] = item;
      });
      other.forEach(function(item, index){
          result[item.vname] = item;
      });
      const out = [];
      Object.keys(result).forEach(function(key) {
          out.push(result[key]);
      });
      return out;
  }

  async newFromFile(contractName, initOptions = [])
  {
    const contractCode = await readContractCodeByName(contractName);

    const init = this.extendInitData([
      // this parameter is mandatory for all init arrays
      {
        vname: '_scilla_version',
        type: 'Uint32',
        value: '0',
      },
      {
        vname: 'owner',
        type: 'ByStr20',
        value: AvelyCore.fromAddress,
      },
    ], initOptions);
    //console.log(init); return;

    const contract = await AvelyCore.zilliqa.contracts.new(contractCode, init);
    return contract;
  }
}

module.exports = new AvelyContracts();
