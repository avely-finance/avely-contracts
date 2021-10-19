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
cdir="runner/$contract"

if [[ ! -d ${cdir} || ! -f ${sdir}/${contract}.scilla || ! -f ${cdir}/state_${i}.json ]]
then
    echo "Test $contract $i does not exist"
    print_usage_and_exit
fi

echo "scilla-runner -init "${cdir}"/init.json -istate "${cdir}/state_${i}".json -imessage "${cdir}/message_${i}".json -o "${cdir}/output_${i}".json -iblockchain "${cdir}/blockchain_${i}".json -i "${sdir}/${contract}".scilla -gaslimit 100000 -libdir /scilla/0/_build/install/default/lib/scilla/stdlib/"

scilla-runner "${cdir}"/init.json -istate "${cdir}/state_${i}".json -imessage "${cdir}/message_${i}".json -o "${cdir}/output_${i}".json -iblockchain "${cdir}/blockchain_${i}".json -i "${sdir}/${contract}".scilla -gaslimit 800000 -libdir /scilla/0/_build/install/default/lib/scilla/stdlib/

status=$?

if test $status -eq 0
then
    echo "output.json emitted by interpreter:"
    cat "${cdir}/output_${i}".json
    echo ""
else
    echo "scilla-runner failed"
fi
