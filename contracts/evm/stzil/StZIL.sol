// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity ^0.8.18;

import {ERC20ZRC2Bridge} from "./ERC20ZRC2Bridge.sol";

interface IZRC2Allowances {
    // non standard ERC20 functions
    function increaseAllowance(address,uint256) external returns (bool);
    function decreaseAllowance(address,uint256) external returns (bool);
}

contract StZIL is ERC20ZRC2Bridge, IZRC2Allowances {

    constructor(address zrc2Address) ERC20ZRC2Bridge(zrc2Address, msg.sender) {
    }

    /* this is a not ERC-20 standard function */
    function increaseAllowance(
        address spender,
        uint256 delta
    ) external override fitsInUint128(delta) returns (bool success) {
        _scillaCallAddressUint128(
            _zrc2Address,
            "IncreaseAllowance",
            spender,
            uint128(delta)
        );
        uint128 currentAllowance128 = _scillaReadMapAddressAddressUint128(
            _zrc2Address,
            "allowances",
            msg.sender,
            spender
        );
        emit Approval(msg.sender, spender, currentAllowance128);
        return true;
    }

    /* this is a not ERC-20 standard function */
    function decreaseAllowance(
        address spender,
        uint256 delta
    ) external override fitsInUint128(delta) returns (bool success) {

        uint128 delta128 = uint128(delta);
        _scillaCallAddressUint128(
            _zrc2Address,
            "DecreaseAllowance",
            spender,
            delta128
        );
        uint128 currentAllowance128 = _scillaReadMapAddressAddressUint128(
            _zrc2Address,
            "allowances",
            msg.sender,
            spender
        );
        emit Approval(msg.sender, spender, currentAllowance128);
        return true;
    }

}
