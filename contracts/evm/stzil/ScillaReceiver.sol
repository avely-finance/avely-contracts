// SPDX-License-Identifier: GPL-3.0-or-later

// solhint-disable one-contract-per-file, func-name-mixedcase

pragma solidity ^0.8.18;

import {ERC165} from "../openzeppelin/contracts/utils/introspection/ERC165.sol";

interface ScillaReceiver {
    function handle_scilla_message(string memory, bytes calldata) external payable;
}

contract ContractSupportingScillaReceiver is ERC165, ScillaReceiver {
    function handle_scilla_message(string memory, bytes calldata) external payable override {
        return;
    }

    function supportsInterface(bytes4 interfaceId) public view virtual override returns (bool isSupported) {
        return interfaceId == type(ScillaReceiver).interfaceId || super.supportsInterface(interfaceId);
    }
}
