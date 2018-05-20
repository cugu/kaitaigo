go run kaitai.go -d kst/$1 tests/$1/formats/*
goimports -w kst/$1
mkdir -p src/kst/$1
mv kst/$1/* src/kst/$1
rm -rf kst

rm -rf src/ks
cp -rf runtime src
mv -f src/runtime src/ks

# run tests
export GOPATH=$GOPATH:$PWD

go test -v tests/$1/test/*
#go test -v -bench=. -run=XXX $1/test/*