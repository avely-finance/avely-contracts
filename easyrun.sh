#!/bin/bash

function print_available_contracts
{
    for ff in tests/runner/*
    do
        if [[ -d $ff ]]
        then
            f=$(basename "$ff")
            echo -n "$f|"
        fi
    done
}

function print_usage_and_exit
{
    echo -n "Usage: $0 ["
    print_available_contracts
    echo "] test_number"
    exit 1
}


if [ $# != 2 ]
then
   print_usage_and_exit
fi

contract=$1
i=$2

sdir="contracts"
cdir="tests/runner/$contract"

if [[ ! -d ${cdir} || ! -f ${sdir}/${contract}.scilla || ! -f ${cdir}/state_${i}.json ]]
then
    echo "Test $contract $i does not exist"
    print_usage_and_exit
fi

echo "=== GDB ==="
echo "file ${SCILLA_RTL_HOME}/bin/scilla-runner"
echo "set args -n "${cdir}"/init.json -s "${cdir}/state_${i}".json -m "${cdir}/message_${i}".json -o "${cdir}/output_${i}".json -b "${cdir}/blockchain_${i}".json -i "${sdir}/${contract}".ll -g 1000000"
echo "dir contracts"
echo "=== === ==="

scilla-runner -init "${cdir}"/init.json -istate "${cdir}/state_${i}".json -imessage "${cdir}/message_${i}".json -o "${cdir}/output_${i}".json -iblockchain "${cdir}/blockchain_${i}".json -i "${sdir}/${contract}".scilla -gaslimit 1000000 -libdir $SCILLA_HOME/stdlib/

status=$?

if test $status -eq 0
then
    echo "output.json emitted by interpreter:"
    cat "${cdir}/output_${i}".json
    echo ""
else
    echo "scilla-runner failed"
fi
