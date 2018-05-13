rm -rf src

go run kaitai.go -d test_formats $1/formats/*
goimports -w test_formats

mkdir -p src
mv test_formats src
cp -r runtime src
mv src/runtime src/kgruntime

# run tests
export GOPATH=$GOPATH:$PWD

go test -v $1/test/*