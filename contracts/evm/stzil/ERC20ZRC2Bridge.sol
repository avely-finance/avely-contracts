// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity ^0.8.18;

import {IERC20} from "../openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from "../openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {ScillaBridge} from "./ScillaBridge.sol";

contract ERC20ZRC2Bridge is IERC20, IERC20Metadata, ScillaBridge {
    address internal _zrc2Address;
    address internal _zrc2Owner;
    uint256[48] private __gap;

    error ValueExceedsUint128Max();

    modifier fitsInUint128(uint256 value) {
        if (value > type(uint128).max) {
            revert ValueExceedsUint128Max();
        }
        _;
    }

    constructor(address zrc2Address, address zrc2Owner) {
        _zrc2Address = zrc2Address;
        _zrc2Owner = zrc2Owner;
    }

    function transfer(address to, uint256 tokens) external override fitsInUint128(tokens) returns (bool success) {
        _scillaCallAddressUint128(_zrc2Address, "Transfer", to, uint128(tokens));
        emit Transfer(msg.sender, to, tokens);
        return true;
    }

    function transferFrom(
        address from,
        address to,
        uint256 tokens
    ) external override fitsInUint128(tokens) returns (bool success) {
        _scillaCallAddressAddressUint128(_zrc2Address, "TransferFrom", from, to, uint128(tokens));
        emit Transfer(from, to, tokens);
        return true;
    }

    function approve(
        address spender,
        uint256 newAllowance
    ) external override fitsInUint128(newAllowance) returns (bool success) {
        uint128 currentAllowance128 = _scillaReadMapAddressAddressUint128(
            _zrc2Address,
            "allowances",
            msg.sender,
            spender
        );
        uint128 newAllowance128 = uint128(newAllowance);
        if (currentAllowance128 >= newAllowance128) {
            _scillaCallAddressUint128(
                _zrc2Address,
                "DecreaseAllowance",
                spender,
                currentAllowance128 - newAllowance128
            );
        } else {
            _scillaCallAddressUint128(
                _zrc2Address,
                "IncreaseAllowance",
                spender,
                newAllowance128 - currentAllowance128
            );
        }
        emit Approval(msg.sender, spender, newAllowance);
        return true;
    }

    function totalSupply() external view override returns (uint256 total) {
        return _scillaReadUint128(_zrc2Address, "total_supply");
    }

    function name() external view override returns (string memory name_) {
        return _scillaReadString(_zrc2Address, "name");
    }

    function symbol() external view override returns (string memory symbol_) {
        return _scillaReadString(_zrc2Address, "symbol");
    }

    function decimals() external view override returns (uint8 retVal) {
        uint256 zilliqaDecimals = _scillaReadUint32(_zrc2Address, "decimals");
        return zilliqaDecimals > 255 ? 255 : uint8(zilliqaDecimals);
    }

    function balanceOf(address tokenOwner) external view override returns (uint256 balance) {
        return _scillaReadMapAddressUint128(_zrc2Address, "balances", tokenOwner);
    }

    function allowance(address tokenOwner, address spender) external view override returns (uint256 allow) {
        return _scillaReadMapAddressAddressUint128(_zrc2Address, "allowances", tokenOwner, spender);
    }
}
