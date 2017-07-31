#!/bin/sh
echo

echo "This will create a new folder named 'mypackage', add files from the genrator"
echo "named 'generic' and run tests on it"
echo "using the commands:"
echo -e "\t 'bp add generic name1 -t DamStruct'"
echo -e "\t 'bp add generic name2 -t int'"
echo -e "\t 'go test -v'"
echo "'bp' should be in path"
echo "run 'bp add generic --help' to understand the flags'"
echo
echo "-- Press any key to continue --"

read

rm -rf mypackage
mkdir -p mypackage
cd mypackage

## BP COMMAND
bp add generic name1 -t DamStruct
bp add generic name2 -t int

go test -v 
cd -
