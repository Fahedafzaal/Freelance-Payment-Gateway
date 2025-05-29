// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// EthJobEscrowMetaData contains all meta data concerning the EthJobEscrow contract.
var EthJobEscrowMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_ethUsdPriceFeed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"FEE_PERCENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"Owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelJob\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"convertUsdToEth\",\"inputs\":[{\"name\":\"usdAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getJobDetails\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"freelancer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ethAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isCompleted\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isPaid\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLatestEthUsd\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"jobs\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"freelancer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ethAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isCompleted\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isPaid\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"markJobCompleted\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"postJob\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"freelancer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"usdAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"client\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"JobCancelled\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"client\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"ethAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"JobCompleted\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"JobPosted\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"client\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"freelancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"usdAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"ethAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentReleased\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"freelancer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"ethAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InsufficientEthSent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"JobAlreadyCompleted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"JobNotCancelable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"JobNotCompleted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotJobClient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyClientCanMarkCompleted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PaymentAlreadyReleased\",\"inputs\":[]}]",
}

// EthJobEscrowABI is the input ABI used to generate the binding from.
// Deprecated: Use EthJobEscrowMetaData.ABI instead.
var EthJobEscrowABI = EthJobEscrowMetaData.ABI

// EthJobEscrow is an auto generated Go binding around an Ethereum contract.
type EthJobEscrow struct {
	EthJobEscrowCaller     // Read-only binding to the contract
	EthJobEscrowTransactor // Write-only binding to the contract
	EthJobEscrowFilterer   // Log filterer for contract events
}

// EthJobEscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthJobEscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthJobEscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthJobEscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthJobEscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthJobEscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthJobEscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthJobEscrowSession struct {
	Contract     *EthJobEscrow     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthJobEscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthJobEscrowCallerSession struct {
	Contract *EthJobEscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// EthJobEscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthJobEscrowTransactorSession struct {
	Contract     *EthJobEscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// EthJobEscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthJobEscrowRaw struct {
	Contract *EthJobEscrow // Generic contract binding to access the raw methods on
}

// EthJobEscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthJobEscrowCallerRaw struct {
	Contract *EthJobEscrowCaller // Generic read-only contract binding to access the raw methods on
}

// EthJobEscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthJobEscrowTransactorRaw struct {
	Contract *EthJobEscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthJobEscrow creates a new instance of EthJobEscrow, bound to a specific deployed contract.
func NewEthJobEscrow(address common.Address, backend bind.ContractBackend) (*EthJobEscrow, error) {
	contract, err := bindEthJobEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthJobEscrow{EthJobEscrowCaller: EthJobEscrowCaller{contract: contract}, EthJobEscrowTransactor: EthJobEscrowTransactor{contract: contract}, EthJobEscrowFilterer: EthJobEscrowFilterer{contract: contract}}, nil
}

// NewEthJobEscrowCaller creates a new read-only instance of EthJobEscrow, bound to a specific deployed contract.
func NewEthJobEscrowCaller(address common.Address, caller bind.ContractCaller) (*EthJobEscrowCaller, error) {
	contract, err := bindEthJobEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthJobEscrowCaller{contract: contract}, nil
}

// NewEthJobEscrowTransactor creates a new write-only instance of EthJobEscrow, bound to a specific deployed contract.
func NewEthJobEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*EthJobEscrowTransactor, error) {
	contract, err := bindEthJobEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthJobEscrowTransactor{contract: contract}, nil
}

// NewEthJobEscrowFilterer creates a new log filterer instance of EthJobEscrow, bound to a specific deployed contract.
func NewEthJobEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*EthJobEscrowFilterer, error) {
	contract, err := bindEthJobEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthJobEscrowFilterer{contract: contract}, nil
}

// bindEthJobEscrow binds a generic wrapper to an already deployed contract.
func bindEthJobEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EthJobEscrowMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthJobEscrow *EthJobEscrowRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthJobEscrow.Contract.EthJobEscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthJobEscrow *EthJobEscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthJobEscrow.Contract.EthJobEscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthJobEscrow *EthJobEscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthJobEscrow.Contract.EthJobEscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthJobEscrow *EthJobEscrowCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthJobEscrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthJobEscrow *EthJobEscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthJobEscrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthJobEscrow *EthJobEscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthJobEscrow.Contract.contract.Transact(opts, method, params...)
}

// FEEPERCENT is a free data retrieval call binding the contract method 0xeaf98d23.
//
// Solidity: function FEE_PERCENT() view returns(uint256)
func (_EthJobEscrow *EthJobEscrowCaller) FEEPERCENT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EthJobEscrow.contract.Call(opts, &out, "FEE_PERCENT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEEPERCENT is a free data retrieval call binding the contract method 0xeaf98d23.
//
// Solidity: function FEE_PERCENT() view returns(uint256)
func (_EthJobEscrow *EthJobEscrowSession) FEEPERCENT() (*big.Int, error) {
	return _EthJobEscrow.Contract.FEEPERCENT(&_EthJobEscrow.CallOpts)
}

// FEEPERCENT is a free data retrieval call binding the contract method 0xeaf98d23.
//
// Solidity: function FEE_PERCENT() view returns(uint256)
func (_EthJobEscrow *EthJobEscrowCallerSession) FEEPERCENT() (*big.Int, error) {
	return _EthJobEscrow.Contract.FEEPERCENT(&_EthJobEscrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0xb4a99a4e.
//
// Solidity: function Owner() view returns(address)
func (_EthJobEscrow *EthJobEscrowCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EthJobEscrow.contract.Call(opts, &out, "Owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0xb4a99a4e.
//
// Solidity: function Owner() view returns(address)
func (_EthJobEscrow *EthJobEscrowSession) Owner() (common.Address, error) {
	return _EthJobEscrow.Contract.Owner(&_EthJobEscrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0xb4a99a4e.
//
// Solidity: function Owner() view returns(address)
func (_EthJobEscrow *EthJobEscrowCallerSession) Owner() (common.Address, error) {
	return _EthJobEscrow.Contract.Owner(&_EthJobEscrow.CallOpts)
}

// ConvertUsdToEth is a free data retrieval call binding the contract method 0xa3053e2a.
//
// Solidity: function convertUsdToEth(uint256 usdAmount) view returns(uint256)
func (_EthJobEscrow *EthJobEscrowCaller) ConvertUsdToEth(opts *bind.CallOpts, usdAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _EthJobEscrow.contract.Call(opts, &out, "convertUsdToEth", usdAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConvertUsdToEth is a free data retrieval call binding the contract method 0xa3053e2a.
//
// Solidity: function convertUsdToEth(uint256 usdAmount) view returns(uint256)
func (_EthJobEscrow *EthJobEscrowSession) ConvertUsdToEth(usdAmount *big.Int) (*big.Int, error) {
	return _EthJobEscrow.Contract.ConvertUsdToEth(&_EthJobEscrow.CallOpts, usdAmount)
}

// ConvertUsdToEth is a free data retrieval call binding the contract method 0xa3053e2a.
//
// Solidity: function convertUsdToEth(uint256 usdAmount) view returns(uint256)
func (_EthJobEscrow *EthJobEscrowCallerSession) ConvertUsdToEth(usdAmount *big.Int) (*big.Int, error) {
	return _EthJobEscrow.Contract.ConvertUsdToEth(&_EthJobEscrow.CallOpts, usdAmount)
}

// GetJobDetails is a free data retrieval call binding the contract method 0x4cac35c6.
//
// Solidity: function getJobDetails(uint256 jobId) view returns(address client, address freelancer, uint256 usdAmount, uint256 ethAmount, bool isCompleted, bool isPaid)
func (_EthJobEscrow *EthJobEscrowCaller) GetJobDetails(opts *bind.CallOpts, jobId *big.Int) (struct {
	Client      common.Address
	Freelancer  common.Address
	UsdAmount   *big.Int
	EthAmount   *big.Int
	IsCompleted bool
	IsPaid      bool
}, error) {
	var out []interface{}
	err := _EthJobEscrow.contract.Call(opts, &out, "getJobDetails", jobId)

	outstruct := new(struct {
		Client      common.Address
		Freelancer  common.Address
		UsdAmount   *big.Int
		EthAmount   *big.Int
		IsCompleted bool
		IsPaid      bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Client = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Freelancer = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.UsdAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.EthAmount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.IsCompleted = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.IsPaid = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// GetJobDetails is a free data retrieval call binding the contract method 0x4cac35c6.
//
// Solidity: function getJobDetails(uint256 jobId) view returns(address client, address freelancer, uint256 usdAmount, uint256 ethAmount, bool isCompleted, bool isPaid)
func (_EthJobEscrow *EthJobEscrowSession) GetJobDetails(jobId *big.Int) (struct {
	Client      common.Address
	Freelancer  common.Address
	UsdAmount   *big.Int
	EthAmount   *big.Int
	IsCompleted bool
	IsPaid      bool
}, error) {
	return _EthJobEscrow.Contract.GetJobDetails(&_EthJobEscrow.CallOpts, jobId)
}

// GetJobDetails is a free data retrieval call binding the contract method 0x4cac35c6.
//
// Solidity: function getJobDetails(uint256 jobId) view returns(address client, address freelancer, uint256 usdAmount, uint256 ethAmount, bool isCompleted, bool isPaid)
func (_EthJobEscrow *EthJobEscrowCallerSession) GetJobDetails(jobId *big.Int) (struct {
	Client      common.Address
	Freelancer  common.Address
	UsdAmount   *big.Int
	EthAmount   *big.Int
	IsCompleted bool
	IsPaid      bool
}, error) {
	return _EthJobEscrow.Contract.GetJobDetails(&_EthJobEscrow.CallOpts, jobId)
}

// GetLatestEthUsd is a free data retrieval call binding the contract method 0xe979ba3f.
//
// Solidity: function getLatestEthUsd() view returns(uint256)
func (_EthJobEscrow *EthJobEscrowCaller) GetLatestEthUsd(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EthJobEscrow.contract.Call(opts, &out, "getLatestEthUsd")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLatestEthUsd is a free data retrieval call binding the contract method 0xe979ba3f.
//
// Solidity: function getLatestEthUsd() view returns(uint256)
func (_EthJobEscrow *EthJobEscrowSession) GetLatestEthUsd() (*big.Int, error) {
	return _EthJobEscrow.Contract.GetLatestEthUsd(&_EthJobEscrow.CallOpts)
}

// GetLatestEthUsd is a free data retrieval call binding the contract method 0xe979ba3f.
//
// Solidity: function getLatestEthUsd() view returns(uint256)
func (_EthJobEscrow *EthJobEscrowCallerSession) GetLatestEthUsd() (*big.Int, error) {
	return _EthJobEscrow.Contract.GetLatestEthUsd(&_EthJobEscrow.CallOpts)
}

// Jobs is a free data retrieval call binding the contract method 0x180aedf3.
//
// Solidity: function jobs(uint256 ) view returns(address client, address freelancer, uint256 usdAmount, uint256 ethAmount, bool isCompleted, bool isPaid)
func (_EthJobEscrow *EthJobEscrowCaller) Jobs(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Client      common.Address
	Freelancer  common.Address
	UsdAmount   *big.Int
	EthAmount   *big.Int
	IsCompleted bool
	IsPaid      bool
}, error) {
	var out []interface{}
	err := _EthJobEscrow.contract.Call(opts, &out, "jobs", arg0)

	outstruct := new(struct {
		Client      common.Address
		Freelancer  common.Address
		UsdAmount   *big.Int
		EthAmount   *big.Int
		IsCompleted bool
		IsPaid      bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Client = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Freelancer = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.UsdAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.EthAmount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.IsCompleted = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.IsPaid = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// Jobs is a free data retrieval call binding the contract method 0x180aedf3.
//
// Solidity: function jobs(uint256 ) view returns(address client, address freelancer, uint256 usdAmount, uint256 ethAmount, bool isCompleted, bool isPaid)
func (_EthJobEscrow *EthJobEscrowSession) Jobs(arg0 *big.Int) (struct {
	Client      common.Address
	Freelancer  common.Address
	UsdAmount   *big.Int
	EthAmount   *big.Int
	IsCompleted bool
	IsPaid      bool
}, error) {
	return _EthJobEscrow.Contract.Jobs(&_EthJobEscrow.CallOpts, arg0)
}

// Jobs is a free data retrieval call binding the contract method 0x180aedf3.
//
// Solidity: function jobs(uint256 ) view returns(address client, address freelancer, uint256 usdAmount, uint256 ethAmount, bool isCompleted, bool isPaid)
func (_EthJobEscrow *EthJobEscrowCallerSession) Jobs(arg0 *big.Int) (struct {
	Client      common.Address
	Freelancer  common.Address
	UsdAmount   *big.Int
	EthAmount   *big.Int
	IsCompleted bool
	IsPaid      bool
}, error) {
	return _EthJobEscrow.Contract.Jobs(&_EthJobEscrow.CallOpts, arg0)
}

// CancelJob is a paid mutator transaction binding the contract method 0x1dffa3dc.
//
// Solidity: function cancelJob(uint256 jobId) returns()
func (_EthJobEscrow *EthJobEscrowTransactor) CancelJob(opts *bind.TransactOpts, jobId *big.Int) (*types.Transaction, error) {
	return _EthJobEscrow.contract.Transact(opts, "cancelJob", jobId)
}

// CancelJob is a paid mutator transaction binding the contract method 0x1dffa3dc.
//
// Solidity: function cancelJob(uint256 jobId) returns()
func (_EthJobEscrow *EthJobEscrowSession) CancelJob(jobId *big.Int) (*types.Transaction, error) {
	return _EthJobEscrow.Contract.CancelJob(&_EthJobEscrow.TransactOpts, jobId)
}

// CancelJob is a paid mutator transaction binding the contract method 0x1dffa3dc.
//
// Solidity: function cancelJob(uint256 jobId) returns()
func (_EthJobEscrow *EthJobEscrowTransactorSession) CancelJob(jobId *big.Int) (*types.Transaction, error) {
	return _EthJobEscrow.Contract.CancelJob(&_EthJobEscrow.TransactOpts, jobId)
}

// MarkJobCompleted is a paid mutator transaction binding the contract method 0x5c1615f3.
//
// Solidity: function markJobCompleted(uint256 jobId) returns()
func (_EthJobEscrow *EthJobEscrowTransactor) MarkJobCompleted(opts *bind.TransactOpts, jobId *big.Int) (*types.Transaction, error) {
	return _EthJobEscrow.contract.Transact(opts, "markJobCompleted", jobId)
}

// MarkJobCompleted is a paid mutator transaction binding the contract method 0x5c1615f3.
//
// Solidity: function markJobCompleted(uint256 jobId) returns()
func (_EthJobEscrow *EthJobEscrowSession) MarkJobCompleted(jobId *big.Int) (*types.Transaction, error) {
	return _EthJobEscrow.Contract.MarkJobCompleted(&_EthJobEscrow.TransactOpts, jobId)
}

// MarkJobCompleted is a paid mutator transaction binding the contract method 0x5c1615f3.
//
// Solidity: function markJobCompleted(uint256 jobId) returns()
func (_EthJobEscrow *EthJobEscrowTransactorSession) MarkJobCompleted(jobId *big.Int) (*types.Transaction, error) {
	return _EthJobEscrow.Contract.MarkJobCompleted(&_EthJobEscrow.TransactOpts, jobId)
}

// PostJob is a paid mutator transaction binding the contract method 0x1892d508.
//
// Solidity: function postJob(uint256 jobId, address freelancer, uint256 usdAmount, address client) payable returns()
func (_EthJobEscrow *EthJobEscrowTransactor) PostJob(opts *bind.TransactOpts, jobId *big.Int, freelancer common.Address, usdAmount *big.Int, client common.Address) (*types.Transaction, error) {
	return _EthJobEscrow.contract.Transact(opts, "postJob", jobId, freelancer, usdAmount, client)
}

// PostJob is a paid mutator transaction binding the contract method 0x1892d508.
//
// Solidity: function postJob(uint256 jobId, address freelancer, uint256 usdAmount, address client) payable returns()
func (_EthJobEscrow *EthJobEscrowSession) PostJob(jobId *big.Int, freelancer common.Address, usdAmount *big.Int, client common.Address) (*types.Transaction, error) {
	return _EthJobEscrow.Contract.PostJob(&_EthJobEscrow.TransactOpts, jobId, freelancer, usdAmount, client)
}

// PostJob is a paid mutator transaction binding the contract method 0x1892d508.
//
// Solidity: function postJob(uint256 jobId, address freelancer, uint256 usdAmount, address client) payable returns()
func (_EthJobEscrow *EthJobEscrowTransactorSession) PostJob(jobId *big.Int, freelancer common.Address, usdAmount *big.Int, client common.Address) (*types.Transaction, error) {
	return _EthJobEscrow.Contract.PostJob(&_EthJobEscrow.TransactOpts, jobId, freelancer, usdAmount, client)
}

// EthJobEscrowJobCancelledIterator is returned from FilterJobCancelled and is used to iterate over the raw logs and unpacked data for JobCancelled events raised by the EthJobEscrow contract.
type EthJobEscrowJobCancelledIterator struct {
	Event *EthJobEscrowJobCancelled // Event containing the contract specifics and raw log

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
func (it *EthJobEscrowJobCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthJobEscrowJobCancelled)
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
		it.Event = new(EthJobEscrowJobCancelled)
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
func (it *EthJobEscrowJobCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthJobEscrowJobCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthJobEscrowJobCancelled represents a JobCancelled event raised by the EthJobEscrow contract.
type EthJobEscrowJobCancelled struct {
	JobId     *big.Int
	Client    common.Address
	EthAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterJobCancelled is a free log retrieval operation binding the contract event 0xa80c76c474b34cc7af71dec63d733b959fff08f4eb0789e288be5db6b608f942.
//
// Solidity: event JobCancelled(uint256 jobId, address indexed client, uint256 ethAmount)
func (_EthJobEscrow *EthJobEscrowFilterer) FilterJobCancelled(opts *bind.FilterOpts, client []common.Address) (*EthJobEscrowJobCancelledIterator, error) {

	var clientRule []interface{}
	for _, clientItem := range client {
		clientRule = append(clientRule, clientItem)
	}

	logs, sub, err := _EthJobEscrow.contract.FilterLogs(opts, "JobCancelled", clientRule)
	if err != nil {
		return nil, err
	}
	return &EthJobEscrowJobCancelledIterator{contract: _EthJobEscrow.contract, event: "JobCancelled", logs: logs, sub: sub}, nil
}

// WatchJobCancelled is a free log subscription operation binding the contract event 0xa80c76c474b34cc7af71dec63d733b959fff08f4eb0789e288be5db6b608f942.
//
// Solidity: event JobCancelled(uint256 jobId, address indexed client, uint256 ethAmount)
func (_EthJobEscrow *EthJobEscrowFilterer) WatchJobCancelled(opts *bind.WatchOpts, sink chan<- *EthJobEscrowJobCancelled, client []common.Address) (event.Subscription, error) {

	var clientRule []interface{}
	for _, clientItem := range client {
		clientRule = append(clientRule, clientItem)
	}

	logs, sub, err := _EthJobEscrow.contract.WatchLogs(opts, "JobCancelled", clientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthJobEscrowJobCancelled)
				if err := _EthJobEscrow.contract.UnpackLog(event, "JobCancelled", log); err != nil {
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

// ParseJobCancelled is a log parse operation binding the contract event 0xa80c76c474b34cc7af71dec63d733b959fff08f4eb0789e288be5db6b608f942.
//
// Solidity: event JobCancelled(uint256 jobId, address indexed client, uint256 ethAmount)
func (_EthJobEscrow *EthJobEscrowFilterer) ParseJobCancelled(log types.Log) (*EthJobEscrowJobCancelled, error) {
	event := new(EthJobEscrowJobCancelled)
	if err := _EthJobEscrow.contract.UnpackLog(event, "JobCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthJobEscrowJobCompletedIterator is returned from FilterJobCompleted and is used to iterate over the raw logs and unpacked data for JobCompleted events raised by the EthJobEscrow contract.
type EthJobEscrowJobCompletedIterator struct {
	Event *EthJobEscrowJobCompleted // Event containing the contract specifics and raw log

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
func (it *EthJobEscrowJobCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthJobEscrowJobCompleted)
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
		it.Event = new(EthJobEscrowJobCompleted)
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
func (it *EthJobEscrowJobCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthJobEscrowJobCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthJobEscrowJobCompleted represents a JobCompleted event raised by the EthJobEscrow contract.
type EthJobEscrowJobCompleted struct {
	JobId *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterJobCompleted is a free log retrieval operation binding the contract event 0x02244c8529cb95e213ee542e76e7776342b3dabd10203d01472bbf4441be8929.
//
// Solidity: event JobCompleted(uint256 jobId)
func (_EthJobEscrow *EthJobEscrowFilterer) FilterJobCompleted(opts *bind.FilterOpts) (*EthJobEscrowJobCompletedIterator, error) {

	logs, sub, err := _EthJobEscrow.contract.FilterLogs(opts, "JobCompleted")
	if err != nil {
		return nil, err
	}
	return &EthJobEscrowJobCompletedIterator{contract: _EthJobEscrow.contract, event: "JobCompleted", logs: logs, sub: sub}, nil
}

// WatchJobCompleted is a free log subscription operation binding the contract event 0x02244c8529cb95e213ee542e76e7776342b3dabd10203d01472bbf4441be8929.
//
// Solidity: event JobCompleted(uint256 jobId)
func (_EthJobEscrow *EthJobEscrowFilterer) WatchJobCompleted(opts *bind.WatchOpts, sink chan<- *EthJobEscrowJobCompleted) (event.Subscription, error) {

	logs, sub, err := _EthJobEscrow.contract.WatchLogs(opts, "JobCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthJobEscrowJobCompleted)
				if err := _EthJobEscrow.contract.UnpackLog(event, "JobCompleted", log); err != nil {
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

// ParseJobCompleted is a log parse operation binding the contract event 0x02244c8529cb95e213ee542e76e7776342b3dabd10203d01472bbf4441be8929.
//
// Solidity: event JobCompleted(uint256 jobId)
func (_EthJobEscrow *EthJobEscrowFilterer) ParseJobCompleted(log types.Log) (*EthJobEscrowJobCompleted, error) {
	event := new(EthJobEscrowJobCompleted)
	if err := _EthJobEscrow.contract.UnpackLog(event, "JobCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthJobEscrowJobPostedIterator is returned from FilterJobPosted and is used to iterate over the raw logs and unpacked data for JobPosted events raised by the EthJobEscrow contract.
type EthJobEscrowJobPostedIterator struct {
	Event *EthJobEscrowJobPosted // Event containing the contract specifics and raw log

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
func (it *EthJobEscrowJobPostedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthJobEscrowJobPosted)
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
		it.Event = new(EthJobEscrowJobPosted)
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
func (it *EthJobEscrowJobPostedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthJobEscrowJobPostedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthJobEscrowJobPosted represents a JobPosted event raised by the EthJobEscrow contract.
type EthJobEscrowJobPosted struct {
	JobId      *big.Int
	Client     common.Address
	Freelancer common.Address
	UsdAmount  *big.Int
	EthAmount  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterJobPosted is a free log retrieval operation binding the contract event 0x5e290dc8e03f46c8c01d0993ae4aff0b815ca77ca9385cf0a5358f66347c04ce.
//
// Solidity: event JobPosted(uint256 jobId, address indexed client, address indexed freelancer, uint256 usdAmount, uint256 ethAmount)
func (_EthJobEscrow *EthJobEscrowFilterer) FilterJobPosted(opts *bind.FilterOpts, client []common.Address, freelancer []common.Address) (*EthJobEscrowJobPostedIterator, error) {

	var clientRule []interface{}
	for _, clientItem := range client {
		clientRule = append(clientRule, clientItem)
	}
	var freelancerRule []interface{}
	for _, freelancerItem := range freelancer {
		freelancerRule = append(freelancerRule, freelancerItem)
	}

	logs, sub, err := _EthJobEscrow.contract.FilterLogs(opts, "JobPosted", clientRule, freelancerRule)
	if err != nil {
		return nil, err
	}
	return &EthJobEscrowJobPostedIterator{contract: _EthJobEscrow.contract, event: "JobPosted", logs: logs, sub: sub}, nil
}

// WatchJobPosted is a free log subscription operation binding the contract event 0x5e290dc8e03f46c8c01d0993ae4aff0b815ca77ca9385cf0a5358f66347c04ce.
//
// Solidity: event JobPosted(uint256 jobId, address indexed client, address indexed freelancer, uint256 usdAmount, uint256 ethAmount)
func (_EthJobEscrow *EthJobEscrowFilterer) WatchJobPosted(opts *bind.WatchOpts, sink chan<- *EthJobEscrowJobPosted, client []common.Address, freelancer []common.Address) (event.Subscription, error) {

	var clientRule []interface{}
	for _, clientItem := range client {
		clientRule = append(clientRule, clientItem)
	}
	var freelancerRule []interface{}
	for _, freelancerItem := range freelancer {
		freelancerRule = append(freelancerRule, freelancerItem)
	}

	logs, sub, err := _EthJobEscrow.contract.WatchLogs(opts, "JobPosted", clientRule, freelancerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthJobEscrowJobPosted)
				if err := _EthJobEscrow.contract.UnpackLog(event, "JobPosted", log); err != nil {
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

// ParseJobPosted is a log parse operation binding the contract event 0x5e290dc8e03f46c8c01d0993ae4aff0b815ca77ca9385cf0a5358f66347c04ce.
//
// Solidity: event JobPosted(uint256 jobId, address indexed client, address indexed freelancer, uint256 usdAmount, uint256 ethAmount)
func (_EthJobEscrow *EthJobEscrowFilterer) ParseJobPosted(log types.Log) (*EthJobEscrowJobPosted, error) {
	event := new(EthJobEscrowJobPosted)
	if err := _EthJobEscrow.contract.UnpackLog(event, "JobPosted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthJobEscrowPaymentReleasedIterator is returned from FilterPaymentReleased and is used to iterate over the raw logs and unpacked data for PaymentReleased events raised by the EthJobEscrow contract.
type EthJobEscrowPaymentReleasedIterator struct {
	Event *EthJobEscrowPaymentReleased // Event containing the contract specifics and raw log

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
func (it *EthJobEscrowPaymentReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthJobEscrowPaymentReleased)
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
		it.Event = new(EthJobEscrowPaymentReleased)
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
func (it *EthJobEscrowPaymentReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthJobEscrowPaymentReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthJobEscrowPaymentReleased represents a PaymentReleased event raised by the EthJobEscrow contract.
type EthJobEscrowPaymentReleased struct {
	JobId      *big.Int
	Freelancer common.Address
	EthAmount  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPaymentReleased is a free log retrieval operation binding the contract event 0x21d71db5be59bb9fa133895586b7404307dd33fb93b16db09dc6f1d9d7d231b0.
//
// Solidity: event PaymentReleased(uint256 jobId, address indexed freelancer, uint256 ethAmount)
func (_EthJobEscrow *EthJobEscrowFilterer) FilterPaymentReleased(opts *bind.FilterOpts, freelancer []common.Address) (*EthJobEscrowPaymentReleasedIterator, error) {

	var freelancerRule []interface{}
	for _, freelancerItem := range freelancer {
		freelancerRule = append(freelancerRule, freelancerItem)
	}

	logs, sub, err := _EthJobEscrow.contract.FilterLogs(opts, "PaymentReleased", freelancerRule)
	if err != nil {
		return nil, err
	}
	return &EthJobEscrowPaymentReleasedIterator{contract: _EthJobEscrow.contract, event: "PaymentReleased", logs: logs, sub: sub}, nil
}

// WatchPaymentReleased is a free log subscription operation binding the contract event 0x21d71db5be59bb9fa133895586b7404307dd33fb93b16db09dc6f1d9d7d231b0.
//
// Solidity: event PaymentReleased(uint256 jobId, address indexed freelancer, uint256 ethAmount)
func (_EthJobEscrow *EthJobEscrowFilterer) WatchPaymentReleased(opts *bind.WatchOpts, sink chan<- *EthJobEscrowPaymentReleased, freelancer []common.Address) (event.Subscription, error) {

	var freelancerRule []interface{}
	for _, freelancerItem := range freelancer {
		freelancerRule = append(freelancerRule, freelancerItem)
	}

	logs, sub, err := _EthJobEscrow.contract.WatchLogs(opts, "PaymentReleased", freelancerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthJobEscrowPaymentReleased)
				if err := _EthJobEscrow.contract.UnpackLog(event, "PaymentReleased", log); err != nil {
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

// ParsePaymentReleased is a log parse operation binding the contract event 0x21d71db5be59bb9fa133895586b7404307dd33fb93b16db09dc6f1d9d7d231b0.
//
// Solidity: event PaymentReleased(uint256 jobId, address indexed freelancer, uint256 ethAmount)
func (_EthJobEscrow *EthJobEscrowFilterer) ParsePaymentReleased(log types.Log) (*EthJobEscrowPaymentReleased, error) {
	event := new(EthJobEscrowPaymentReleased)
	if err := _EthJobEscrow.contract.UnpackLog(event, "PaymentReleased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
