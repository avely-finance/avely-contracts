// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity ^0.8.18;

import {ERC20ZRC2Bridge} from "./ERC20ZRC2Bridge.sol";

contract StZIL is ERC20ZRC2Bridge {
    constructor(address zrc2Address) ERC20ZRC2Bridge(zrc2Address, msg.sender) {
    }
}
