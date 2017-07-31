rm -rf mypackage
mkdir -p mypackage
cd mypackage
echo "Running bp command"
go run ../../../cli/bp/main.go add generic name1 -t DamStruct
go run ../../../cli/bp/main.go add generic name2 -t int
go test -v 
cd -
