#!/bin/bash


contract=$1

sdir="contracts"
cdir="tests/runner/$contract"

if [[ ! -f ${sdir}/${contract}.scilla ]]
then
    echo "Contract $contract does not exist"
    exit 1
fi

scilla-checker -init "${cdir}"/init.json -gaslimit 1000000 -libdir $SCILLA_HOME/stdlib/ "${sdir}/${contract}".scilla
