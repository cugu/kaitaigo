# install
echo install
goimports -w cmd
go install gitlab.com/cugu/kaitai.go/cmd/...

# generate code
echo generate code
compiler parser/...

# test
echo go test
go test -v `go list gitlab.com/cugu/kaitai.go/parser/... 2> /dev/null | grep -v "/kaitai/"`

# go test --bench=. -v gitlab.com/cugu/kaitai.go/parser/...