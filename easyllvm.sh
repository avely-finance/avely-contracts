sdir="contracts"
contract=$1

if [[ ! -f ${sdir}/${contract}.scilla ]]
then
    echo "Test $contract does not exist"
fi

scilla-llvm -libdir $SCILLA_HOME/stdlib -gaslimit 100000 -debuginfo true contracts/$contract.scilla > contracts/$contract.ll

echo "OK"
