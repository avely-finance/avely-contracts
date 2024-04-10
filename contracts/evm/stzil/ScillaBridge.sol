// SPDX-License-Identifier: GPL-3.0-or-later

// solhint-disable no-inline-assembly

pragma solidity ^0.8.18;

contract ScillaBridge {
    // see https://github.com/Zilliqa/ZIP/blob/7e7e1b3cb1b0faa612427af5e8af3b174207657b/zips/zip-21.md#scilla_call_2
    uint8 private constant _CALL_MODE_SENDER_IS_CONTRACT = 0;
    uint8 private constant _CALL_MODE_SENDER_IS_MSG_SENDER = 1;
    uint256 private constant _CALL_GAS = 21000;
    uint256 private constant _STATIC_CALL_GAS = 21000;
    address private constant _SCILLA_READ_PRECOMPILE = 0x000000000000000000000000000000005A494C92;
    address private constant _SCILLA_CALL2_PRECOMPILE = 0x000000000000000000000000000000005a494c53;

    error ScillaCallFailed();
    error ScillaStaticCallFailed();

    function _scillaCallWithValue(address contractAddress, string memory transitionName) internal {
        bytes memory encodedArgs = abi.encode(contractAddress, transitionName, _CALL_MODE_SENDER_IS_MSG_SENDER);
        uint256 argsLength = encodedArgs.length;
        bool success;
        assembly {
            success := call(
                _CALL_GAS,
                _SCILLA_CALL2_PRECOMPILE,
                callvalue(),
                add(encodedArgs, 0x20),
                argsLength,
                0x20,
                0
            )
        }
        if (!success) {
            revert ScillaCallFailed();
        }
    }

    function _scillaCall(address contractAddress, string memory transitionName) internal {
        bytes memory encodedArgs = abi.encode(contractAddress, transitionName, _CALL_MODE_SENDER_IS_MSG_SENDER);
        uint256 argsLength = encodedArgs.length;
        bool success;
        assembly {
            success := call(_CALL_GAS, _SCILLA_CALL2_PRECOMPILE, 0, add(encodedArgs, 0x20), argsLength, 0x20, 0)
        }
        if (!success) {
            revert ScillaCallFailed();
        }
    }

    function _scillaCallUint128(
        address contractAddress,
        string memory transitionName,
        uint128 amount
    ) internal {
        bytes memory encodedArgs = abi.encode(
            contractAddress,
            transitionName,
            _CALL_MODE_SENDER_IS_MSG_SENDER,
            amount
        );
        uint256 argsLength = encodedArgs.length;
        bool success;
        assembly {
            success := call(_CALL_GAS, _SCILLA_CALL2_PRECOMPILE, 0, add(encodedArgs, 0x20), argsLength, 0x20, 0)
        }
        if (!success) {
            revert ScillaCallFailed();
        }
    }

    function _scillaCallAddress(address contractAddress, string memory transitionName, address addr1) internal {
        bytes memory encodedArgs = abi.encode(contractAddress, transitionName, _CALL_MODE_SENDER_IS_MSG_SENDER, addr1);
        uint256 argsLength = encodedArgs.length;
        bool success;
        assembly {
            success := call(_CALL_GAS, _SCILLA_CALL2_PRECOMPILE, 0, add(encodedArgs, 0x20), argsLength, 0x20, 0)
        }
        if (!success) {
            revert ScillaCallFailed();
        }
    }

    function _scillaCallAddressUint128(
        address contractAddress,
        string memory transitionName,
        address addr1,
        uint128 amount
    ) internal {
        bytes memory encodedArgs = abi.encode(
            contractAddress,
            transitionName,
            _CALL_MODE_SENDER_IS_MSG_SENDER,
            addr1,
            amount
        );
        uint256 argsLength = encodedArgs.length;
        bool success;
        assembly {
            success := call(_CALL_GAS, _SCILLA_CALL2_PRECOMPILE, 0, add(encodedArgs, 0x20), argsLength, 0x20, 0)
        }
        if (!success) {
            revert ScillaCallFailed();
        }
    }

    function _scillaCallAddressAddressUint128(
        address contractAddress,
        string memory transitionName,
        address addr1,
        address addr2,
        uint128 amount
    ) internal {
        // solhint-disable-next-line func-named-parameters
        bytes memory encodedArgs = abi.encode(
            contractAddress,
            transitionName,
            _CALL_MODE_SENDER_IS_MSG_SENDER,
            addr1,
            addr2,
            amount
        );
        uint256 argsLength = encodedArgs.length;
        bool success;
        assembly {
            success := call(_CALL_GAS, _SCILLA_CALL2_PRECOMPILE, 0, add(encodedArgs, 0x20), argsLength, 0x20, 0)
        }
        if (!success) {
            revert ScillaCallFailed();
        }
    }

    function _scillaReadUint128(address contractAddress, string memory varName) internal view returns (uint128 retVal) {
        bytes memory encodedArgs = abi.encode(contractAddress, varName);
        uint256 argsLength = encodedArgs.length;
        bool success;
        bytes memory output = new bytes(36);
        assembly {
            success := staticcall(
                _STATIC_CALL_GAS,
                _SCILLA_READ_PRECOMPILE,
                add(encodedArgs, 0x20),
                argsLength,
                add(output, 0x20),
                32
            )
        }
        if (!success) {
            revert ScillaStaticCallFailed();
        }
        retVal = abi.decode(output, (uint128));
    }

    function _scillaReadUint32(address contractAddress, string memory varName) internal view returns (uint32 retVal) {
        bytes memory encodedArgs = abi.encode(contractAddress, varName);
        uint256 argsLength = encodedArgs.length;
        bool success;
        bytes memory output = new bytes(36);
        assembly {
            success := staticcall(
                _STATIC_CALL_GAS,
                _SCILLA_READ_PRECOMPILE,
                add(encodedArgs, 0x20),
                argsLength,
                add(output, 0x20),
                32
            )
        }
        if (!success) {
            revert ScillaStaticCallFailed();
        }
        retVal = abi.decode(output, (uint32));
    }

    function _scillaReadMapAddressUint128(
        address contractAddress,
        string memory varName,
        address key
    ) internal view returns (uint128 retVal) {
        bytes memory encodedArgs = abi.encode(contractAddress, varName, key);
        uint256 argsLength = encodedArgs.length;
        bool success;
        bytes memory output = new bytes(36);
        assembly {
            success := staticcall(
                _STATIC_CALL_GAS,
                _SCILLA_READ_PRECOMPILE,
                add(encodedArgs, 0x20),
                argsLength,
                add(output, 0x20),
                32
            )
        }
        if (!success) {
            revert ScillaStaticCallFailed();
        }
        retVal = abi.decode(output, (uint128));
    }

    function _scillaReadMapAddressAddressUint128(
        address contractAddress,
        string memory varName,
        address key1,
        address key2
    ) internal view returns (uint128 retVal) {
        bytes memory encodedArgs = abi.encode(contractAddress, varName, key1, key2);
        uint256 argsLength = encodedArgs.length;
        bool success;
        bytes memory output = new bytes(36);
        assembly {
            success := staticcall(
                _STATIC_CALL_GAS,
                _SCILLA_READ_PRECOMPILE,
                add(encodedArgs, 0x20),
                argsLength,
                add(output, 0x20),
                32
            )
        }
        if (!success) {
            revert ScillaStaticCallFailed();
        }
        retVal = abi.decode(output, (uint128));
    }

    function _scillaReadString(
        address contractAddress,
        string memory varName
    ) internal view returns (string memory retVal) {
        bytes memory encodedArgs = abi.encode(contractAddress, varName);
        uint256 argsLength = encodedArgs.length;
        bool success;
        bytes memory output = new bytes(256);
        uint256 outputLen = output.length - 4;
        assembly {
            success := staticcall(
                _STATIC_CALL_GAS,
                _SCILLA_READ_PRECOMPILE,
                add(encodedArgs, 0x20),
                argsLength,
                add(output, 0x20),
                outputLen
            )
        }
        if (!success) {
            revert ScillaStaticCallFailed();
        }

        retVal = abi.decode(output, (string));
    }
}
