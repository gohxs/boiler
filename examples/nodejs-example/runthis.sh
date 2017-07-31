#!/bin/sh
echo
echo "This will create a new project named 'proj' based on nodejs-boilerplate"
echo "using the command: 'bp new ./nodejs-boilerplate proj'"
echo "'bp' should be in path"
echo
echo "-- Press any key to continue --"

read

bp new nodejs-boilerplate proj

cat proj/package.json
