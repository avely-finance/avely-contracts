# Debug Scilla

Articles:
- [Debug Scilla](https://blog.zilliqa.com/a-debugger-for-scilla-the-beginning-97eafd562206)
- [Scilla Compiler](https://blog.zilliqa.com/a-compiled-backend-for-scilla-38e9a34630e0)

## Requirements

_Note: You don't need to install additional tools if you use the current project's Docker setup_

1. [Scilla](https://github.com/Zilliqa/scilla)
2. [LLVM with OCaml bindings](https://github.com/Zilliqa/scilla-compiler/blob/master/scripts/build_install_llvm.sh)
3. [Scilla Compiler](https://github.com/Zilliqa/scilla-compiler)
4. [Scilla RTL](https://github.com/Zilliqa/scilla-rtl)
5. [GDB](https://www.gnu.org/software/gdb/)

## Steps

First of all, we need to compile a contract to LLVM IR:

```sh
$ scilla-llvm -libdir $SCILLA_HOME/stdlib -gaslimit 100000 -debuginfo true contracts/<YOUR_CONTRACT>.scilla > contracts/<YOUR_CONTRACT>.ll
```

To debug a contract execution we need to prepare a test suite. Check out` scilla-runner --help` for the instructions.
We have a simple wrapper to easily run `scilla-runner` with a convention used in the official repo:

```sh
$ ./easyrun.sh <YOUR_CONTRACT_NAME> <TEST_SUITE_NUMBER>
```

This command returns similar output:

```sh
=== GDB ===
🐞 file /root/tools/scilla-rtl/build/bin/scilla-runner
🐞 set args -n tests/runner/stubStakingContract/init.json -s tests/runner/stubStakingContract/state_1.json -m tests/runner/stubStakingContract/message_1.json -o tests/runner/stubStakingContract/output_1.json -b tests/runner/stubStakingContract/blockchain_1.json -i contracts/stubStakingContract.ll -g 1000000
🐞 dir contracts
=== === ===
...
```


We have everything we need to run a debug session:

```sh
$ gdb

<PUT HERE COMMANDS FROM easyrun output>

b <nameOfYourFunction>

run
```

[See a demo for more information](https://www.youtube.com/watch?v=N3T5cp9cbvU)

Enjoy!
