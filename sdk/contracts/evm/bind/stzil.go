// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bind

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// StZILMetaData contains all meta data concerning the StZIL contract.
var StZILMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"zrc2Address\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ScillaCallFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ScillaStaticCallFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValueExceedsUint128Max\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenOwner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"allow\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newAllowance\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenOwner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"retVal\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"delta\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"foo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"delta\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"increaseAllowance1\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600860325534801561001557600080fd5b5060405161102338038061102383398101604081905261003491610062565b600080546001600160a01b039092166001600160a01b03199283161790556001805490911633179055610092565b60006020828403121561007457600080fd5b81516001600160a01b038116811461008b57600080fd5b9392505050565b610f82806100a16000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c80635b3cc7711161008c578063a457c2d711610066578063a457c2d714610192578063a9059cbb146101a5578063c2985578146101b8578063dd62ed3e146101c157600080fd5b80635b3cc7711461016b57806370a082311461017757806395d89b411461018a57600080fd5b806306fdde03146100d4578063095ea7b3146100f257806318160ddd1461011557806323b872dd1461012b578063313ce5671461013e5780633950935114610158575b600080fd5b6100dc6101d4565b6040516100e99190610bd1565b60405180910390f35b610105610100366004610c07565b61020f565b60405190151581526020016100e9565b61011d61036b565b6040519081526020016100e9565b610105610139366004610c31565b6103b3565b610146610472565b60405160ff90911681526020016100e9565b610105610166366004610c07565b6104cb565b60016032819055610105565b61011d610185366004610c6d565b61051b565b6100dc610561565b6101056101a0366004610c07565b610599565b6101056101b3366004610c07565b61068b565b61011d60325481565b61011d6101cf366004610c88565b610737565b6000546040805180820190915260048152636e616d6560e01b602082015260609161020a916001600160a01b0390911690610781565b905090565b6000816001600160801b0381111561023a5760405163568d99a160e01b815260040160405180910390fd5b6000805460408051808201909152600a815269616c6c6f77616e63657360b01b6020820152610274916001600160a01b0316903388610836565b9050836001600160801b03808216908316106102d8576000546040805180820190915260118152704465637265617365416c6c6f77616e636560781b60208201526102d3916001600160a01b031690886102ce8587610cd1565b6108db565b61031c565b600054604080518082019091526011815270496e637265617365416c6c6f77616e636560781b602082015261031c916001600160a01b031690886102ce8686610cd1565b6040518581526001600160a01b0387169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925906020015b60405180910390a36001935050505b5092915050565b6000805460408051808201909152600c81526b746f74616c5f737570706c7960a01b60208201526103a5916001600160a01b031690610948565b6001600160801b0316905090565b6000816001600160801b038111156103de5760405163568d99a160e01b815260040160405180910390fd5b60005460408051808201909152600c81526b5472616e7366657246726f6d60a01b602082015261041a916001600160a01b0316908787876109e7565b836001600160a01b0316856001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8560405161045f91815260200190565b60405180910390a3506001949350505050565b60008054604080518082019091526008815267646563696d616c7360c01b602082015282916104ac916001600160a01b0390911690610a57565b63ffffffff16905060ff81116104c257806104c5565b60ff5b91505090565b600260325560008054604080518082019091526011815270496e637265617365416c6c6f77616e636560781b6020820152610511916001600160a01b03169085856108db565b5060015b92915050565b6000805460408051808201909152600881526762616c616e63657360c01b6020820152610552916001600160a01b03169084610aeb565b6001600160801b031692915050565b6000546040805180820190915260068152651cde5b589bdb60d21b602082015260609161020a916001600160a01b0390911690610781565b6000816001600160801b038111156105c45760405163568d99a160e01b815260040160405180910390fd5b6000546040805180820190915260118152704465637265617365416c6c6f77616e636560781b60208201528491610608916001600160a01b039091169087846108db565b6000805460408051808201909152600a815269616c6c6f77616e63657360b01b6020820152610642916001600160a01b0316903389610836565b6040516001600160801b03821681529091506001600160a01b0387169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92590602001610355565b6000816001600160801b038111156106b65760405163568d99a160e01b815260040160405180910390fd5b6000546040805180820190915260088152672a3930b739b332b960c11b60208201526106ed916001600160a01b03169086866108db565b6040518381526001600160a01b0385169033907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35060019392505050565b6000805460408051808201909152600a815269616c6c6f77616e63657360b01b6020820152610771916001600160a01b0316908585610836565b6001600160801b03169392505050565b606060008383604051602001610798929190610cf1565b60408051808303601f1901815282825280516101008085526101208501909352909350916000918291906020820181803683370190505090506000600482516107e19190610d33565b905080602083018560208801635a494c92615208fa9250826108165760405163bd697fed60e01b815260040160405180910390fd5b8180602001905181019061082a9190610d46565b98975050505050505050565b600080858585856040516020016108509493929190610de8565b60408051808303601f1901815282825280516024808552606085019093529093509160009182919060208201818036833701905050905060208082018460208701635a494c92615208fa9150816108ba5760405163bd697fed60e01b815260040160405180910390fd5b808060200190518101906108ce9190610e23565b9998505050505050505050565b60008484600185856040516020016108f7959493929190610e4c565b60408051601f19818403018152919052805190915060008060208381860183635a494c53615208f190508061093f5760405163426b730b60e11b815260040160405180910390fd5b50505050505050565b600080838360405160200161095e929190610cf1565b60408051808303601f1901815282825280516024808552606085019093529093509160009182919060208201818036833701905050905060208082018460208701635a494c92615208fa9150816109c85760405163bd697fed60e01b815260040160405180910390fd5b808060200190518101906109dc9190610e23565b979650505050505050565b600085856001868686604051602001610a0596959493929190610e9b565b60408051601f19818403018152919052805190915060008060208381860183635a494c53615208f1905080610a4d5760405163426b730b60e11b815260040160405180910390fd5b5050505050505050565b6000808383604051602001610a6d929190610cf1565b60408051808303601f1901815282825280516024808552606085019093529093509160009182919060208201818036833701905050905060208082018460208701635a494c92615208fa915081610ad75760405163bd697fed60e01b815260040160405180910390fd5b808060200190518101906109dc9190610ef1565b600080848484604051602001610b0393929190610f17565b60408051808303601f1901815282825280516024808552606085019093529093509160009182919060208201818036833701905050905060208082018460208701635a494c92615208fa915081610b6d5760405163bd697fed60e01b815260040160405180910390fd5b8080602001905181019061082a9190610e23565b60005b83811015610b9c578181015183820152602001610b84565b50506000910152565b60008151808452610bbd816020860160208601610b81565b601f01601f19169290920160200192915050565b602081526000610be46020830184610ba5565b9392505050565b80356001600160a01b0381168114610c0257600080fd5b919050565b60008060408385031215610c1a57600080fd5b610c2383610beb565b946020939093013593505050565b600080600060608486031215610c4657600080fd5b610c4f84610beb565b9250610c5d60208501610beb565b9150604084013590509250925092565b600060208284031215610c7f57600080fd5b610be482610beb565b60008060408385031215610c9b57600080fd5b610ca483610beb565b9150610cb260208401610beb565b90509250929050565b634e487b7160e01b600052601160045260246000fd5b6001600160801b0382811682821603908082111561036457610364610cbb565b6001600160a01b0383168152604060208201819052600090610d1590830184610ba5565b949350505050565b634e487b7160e01b600052604160045260246000fd5b8181038181111561051557610515610cbb565b600060208284031215610d5857600080fd5b815167ffffffffffffffff80821115610d7057600080fd5b818401915084601f830112610d8457600080fd5b815181811115610d9657610d96610d1d565b604051601f8201601f19908116603f01168101908382118183101715610dbe57610dbe610d1d565b81604052828152876020848701011115610dd757600080fd5b6109dc836020830160208801610b81565b600060018060a01b03808716835260806020840152610e0a6080840187610ba5565b9481166040840152929092166060909101525092915050565b600060208284031215610e3557600080fd5b81516001600160801b0381168114610be457600080fd5b600060018060a01b03808816835260a06020840152610e6e60a0840188610ba5565b915060ff861660408401528085166060840152506001600160801b03831660808301529695505050505050565b600060018060a01b03808916835260c06020840152610ebd60c0840189610ba5565b60ff9790971660408401529485166060830152509190921660808201526001600160801b0390911660a09091015292915050565b600060208284031215610f0357600080fd5b815163ffffffff81168114610be457600080fd5b600060018060a01b03808616835260606020840152610f396060840186610ba5565b915080841660408401525094935050505056fea26469706673582212203b60cdab653ce4b5b3668217976fbcb2feb73938bef1fcdd22d275f47a663f0064736f6c63430008130033",
}

// StZILABI is the input ABI used to generate the binding from.
// Deprecated: Use StZILMetaData.ABI instead.
var StZILABI = StZILMetaData.ABI

// StZILBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StZILMetaData.Bin instead.
var StZILBin = StZILMetaData.Bin

// DeployStZIL deploys a new Ethereum contract, binding an instance of StZIL to it.
func DeployStZIL(auth *bind.TransactOpts, backend bind.ContractBackend, zrc2Address common.Address) (common.Address, *types.Transaction, *StZIL, error) {
	parsed, err := StZILMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StZILBin), backend, zrc2Address)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StZIL{StZILCaller: StZILCaller{contract: contract}, StZILTransactor: StZILTransactor{contract: contract}, StZILFilterer: StZILFilterer{contract: contract}}, nil
}

// StZIL is an auto generated Go binding around an Ethereum contract.
type StZIL struct {
	StZILCaller     // Read-only binding to the contract
	StZILTransactor // Write-only binding to the contract
	StZILFilterer   // Log filterer for contract events
}

// StZILCaller is an auto generated read-only Go binding around an Ethereum contract.
type StZILCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StZILTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StZILTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StZILFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StZILFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StZILSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StZILSession struct {
	Contract     *StZIL            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StZILCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StZILCallerSession struct {
	Contract *StZILCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StZILTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StZILTransactorSession struct {
	Contract     *StZILTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StZILRaw is an auto generated low-level Go binding around an Ethereum contract.
type StZILRaw struct {
	Contract *StZIL // Generic contract binding to access the raw methods on
}

// StZILCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StZILCallerRaw struct {
	Contract *StZILCaller // Generic read-only contract binding to access the raw methods on
}

// StZILTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StZILTransactorRaw struct {
	Contract *StZILTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStZIL creates a new instance of StZIL, bound to a specific deployed contract.
func NewStZIL(address common.Address, backend bind.ContractBackend) (*StZIL, error) {
	contract, err := bindStZIL(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StZIL{StZILCaller: StZILCaller{contract: contract}, StZILTransactor: StZILTransactor{contract: contract}, StZILFilterer: StZILFilterer{contract: contract}}, nil
}

// NewStZILCaller creates a new read-only instance of StZIL, bound to a specific deployed contract.
func NewStZILCaller(address common.Address, caller bind.ContractCaller) (*StZILCaller, error) {
	contract, err := bindStZIL(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StZILCaller{contract: contract}, nil
}

// NewStZILTransactor creates a new write-only instance of StZIL, bound to a specific deployed contract.
func NewStZILTransactor(address common.Address, transactor bind.ContractTransactor) (*StZILTransactor, error) {
	contract, err := bindStZIL(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StZILTransactor{contract: contract}, nil
}

// NewStZILFilterer creates a new log filterer instance of StZIL, bound to a specific deployed contract.
func NewStZILFilterer(address common.Address, filterer bind.ContractFilterer) (*StZILFilterer, error) {
	contract, err := bindStZIL(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StZILFilterer{contract: contract}, nil
}

// bindStZIL binds a generic wrapper to an already deployed contract.
func bindStZIL(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StZILMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StZIL *StZILRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StZIL.Contract.StZILCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StZIL *StZILRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StZIL.Contract.StZILTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StZIL *StZILRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StZIL.Contract.StZILTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StZIL *StZILCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StZIL.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StZIL *StZILTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StZIL.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StZIL *StZILTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StZIL.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address tokenOwner, address spender) view returns(uint256 allow)
func (_StZIL *StZILCaller) Allowance(opts *bind.CallOpts, tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StZIL.contract.Call(opts, &out, "allowance", tokenOwner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address tokenOwner, address spender) view returns(uint256 allow)
func (_StZIL *StZILSession) Allowance(tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	return _StZIL.Contract.Allowance(&_StZIL.CallOpts, tokenOwner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address tokenOwner, address spender) view returns(uint256 allow)
func (_StZIL *StZILCallerSession) Allowance(tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	return _StZIL.Contract.Allowance(&_StZIL.CallOpts, tokenOwner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) view returns(uint256 balance)
func (_StZIL *StZILCaller) BalanceOf(opts *bind.CallOpts, tokenOwner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StZIL.contract.Call(opts, &out, "balanceOf", tokenOwner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) view returns(uint256 balance)
func (_StZIL *StZILSession) BalanceOf(tokenOwner common.Address) (*big.Int, error) {
	return _StZIL.Contract.BalanceOf(&_StZIL.CallOpts, tokenOwner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) view returns(uint256 balance)
func (_StZIL *StZILCallerSession) BalanceOf(tokenOwner common.Address) (*big.Int, error) {
	return _StZIL.Contract.BalanceOf(&_StZIL.CallOpts, tokenOwner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 retVal)
func (_StZIL *StZILCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _StZIL.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 retVal)
func (_StZIL *StZILSession) Decimals() (uint8, error) {
	return _StZIL.Contract.Decimals(&_StZIL.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8 retVal)
func (_StZIL *StZILCallerSession) Decimals() (uint8, error) {
	return _StZIL.Contract.Decimals(&_StZIL.CallOpts)
}

// Foo is a free data retrieval call binding the contract method 0xc2985578.
//
// Solidity: function foo() view returns(uint256)
func (_StZIL *StZILCaller) Foo(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StZIL.contract.Call(opts, &out, "foo")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Foo is a free data retrieval call binding the contract method 0xc2985578.
//
// Solidity: function foo() view returns(uint256)
func (_StZIL *StZILSession) Foo() (*big.Int, error) {
	return _StZIL.Contract.Foo(&_StZIL.CallOpts)
}

// Foo is a free data retrieval call binding the contract method 0xc2985578.
//
// Solidity: function foo() view returns(uint256)
func (_StZIL *StZILCallerSession) Foo() (*big.Int, error) {
	return _StZIL.Contract.Foo(&_StZIL.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string name_)
func (_StZIL *StZILCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StZIL.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string name_)
func (_StZIL *StZILSession) Name() (string, error) {
	return _StZIL.Contract.Name(&_StZIL.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string name_)
func (_StZIL *StZILCallerSession) Name() (string, error) {
	return _StZIL.Contract.Name(&_StZIL.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string symbol_)
func (_StZIL *StZILCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StZIL.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string symbol_)
func (_StZIL *StZILSession) Symbol() (string, error) {
	return _StZIL.Contract.Symbol(&_StZIL.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string symbol_)
func (_StZIL *StZILCallerSession) Symbol() (string, error) {
	return _StZIL.Contract.Symbol(&_StZIL.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 total)
func (_StZIL *StZILCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StZIL.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 total)
func (_StZIL *StZILSession) TotalSupply() (*big.Int, error) {
	return _StZIL.Contract.TotalSupply(&_StZIL.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256 total)
func (_StZIL *StZILCallerSession) TotalSupply() (*big.Int, error) {
	return _StZIL.Contract.TotalSupply(&_StZIL.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 newAllowance) returns(bool success)
func (_StZIL *StZILTransactor) Approve(opts *bind.TransactOpts, spender common.Address, newAllowance *big.Int) (*types.Transaction, error) {
	return _StZIL.contract.Transact(opts, "approve", spender, newAllowance)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 newAllowance) returns(bool success)
func (_StZIL *StZILSession) Approve(spender common.Address, newAllowance *big.Int) (*types.Transaction, error) {
	return _StZIL.Contract.Approve(&_StZIL.TransactOpts, spender, newAllowance)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 newAllowance) returns(bool success)
func (_StZIL *StZILTransactorSession) Approve(spender common.Address, newAllowance *big.Int) (*types.Transaction, error) {
	return _StZIL.Contract.Approve(&_StZIL.TransactOpts, spender, newAllowance)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 delta) returns(bool success)
func (_StZIL *StZILTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, delta *big.Int) (*types.Transaction, error) {
	return _StZIL.contract.Transact(opts, "decreaseAllowance", spender, delta)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 delta) returns(bool success)
func (_StZIL *StZILSession) DecreaseAllowance(spender common.Address, delta *big.Int) (*types.Transaction, error) {
	return _StZIL.Contract.DecreaseAllowance(&_StZIL.TransactOpts, spender, delta)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 delta) returns(bool success)
func (_StZIL *StZILTransactorSession) DecreaseAllowance(spender common.Address, delta *big.Int) (*types.Transaction, error) {
	return _StZIL.Contract.DecreaseAllowance(&_StZIL.TransactOpts, spender, delta)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 delta) returns(bool success)
func (_StZIL *StZILTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, delta *big.Int) (*types.Transaction, error) {
	return _StZIL.contract.Transact(opts, "increaseAllowance", spender, delta)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 delta) returns(bool success)
func (_StZIL *StZILSession) IncreaseAllowance(spender common.Address, delta *big.Int) (*types.Transaction, error) {
	return _StZIL.Contract.IncreaseAllowance(&_StZIL.TransactOpts, spender, delta)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 delta) returns(bool success)
func (_StZIL *StZILTransactorSession) IncreaseAllowance(spender common.Address, delta *big.Int) (*types.Transaction, error) {
	return _StZIL.Contract.IncreaseAllowance(&_StZIL.TransactOpts, spender, delta)
}

// IncreaseAllowance1 is a paid mutator transaction binding the contract method 0x5b3cc771.
//
// Solidity: function increaseAllowance1() returns(bool success)
func (_StZIL *StZILTransactor) IncreaseAllowance1(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StZIL.contract.Transact(opts, "increaseAllowance1")
}

// IncreaseAllowance1 is a paid mutator transaction binding the contract method 0x5b3cc771.
//
// Solidity: function increaseAllowance1() returns(bool success)
func (_StZIL *StZILSession) IncreaseAllowance1() (*types.Transaction, error) {
	return _StZIL.Contract.IncreaseAllowance1(&_StZIL.TransactOpts)
}

// IncreaseAllowance1 is a paid mutator transaction binding the contract method 0x5b3cc771.
//
// Solidity: function increaseAllowance1() returns(bool success)
func (_StZIL *StZILTransactorSession) IncreaseAllowance1() (*types.Transaction, error) {
	return _StZIL.Contract.IncreaseAllowance1(&_StZIL.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_StZIL *StZILTransactor) Transfer(opts *bind.TransactOpts, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _StZIL.contract.Transact(opts, "transfer", to, tokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_StZIL *StZILSession) Transfer(to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _StZIL.Contract.Transfer(&_StZIL.TransactOpts, to, tokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 tokens) returns(bool success)
func (_StZIL *StZILTransactorSession) Transfer(to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _StZIL.Contract.Transfer(&_StZIL.TransactOpts, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_StZIL *StZILTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _StZIL.contract.Transact(opts, "transferFrom", from, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_StZIL *StZILSession) TransferFrom(from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _StZIL.Contract.TransferFrom(&_StZIL.TransactOpts, from, to, tokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokens) returns(bool success)
func (_StZIL *StZILTransactorSession) TransferFrom(from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return _StZIL.Contract.TransferFrom(&_StZIL.TransactOpts, from, to, tokens)
}

// StZILApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the StZIL contract.
type StZILApprovalIterator struct {
	Event *StZILApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StZILApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StZILApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StZILApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StZILApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StZILApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StZILApproval represents a Approval event raised by the StZIL contract.
type StZILApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_StZIL *StZILFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*StZILApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StZIL.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &StZILApprovalIterator{contract: _StZIL.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_StZIL *StZILFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *StZILApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _StZIL.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StZILApproval)
				if err := _StZIL.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_StZIL *StZILFilterer) ParseApproval(log types.Log) (*StZILApproval, error) {
	event := new(StZILApproval)
	if err := _StZIL.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StZILTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the StZIL contract.
type StZILTransferIterator struct {
	Event *StZILTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StZILTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StZILTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StZILTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StZILTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StZILTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StZILTransfer represents a Transfer event raised by the StZIL contract.
type StZILTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_StZIL *StZILFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*StZILTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StZIL.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &StZILTransferIterator{contract: _StZIL.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_StZIL *StZILFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StZILTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _StZIL.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StZILTransfer)
				if err := _StZIL.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_StZIL *StZILFilterer) ParseTransfer(log types.Log) (*StZILTransfer, error) {
	event := new(StZILTransfer)
	if err := _StZIL.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
