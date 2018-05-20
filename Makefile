all: clean install code test

clean:
	find . -name "*.ksy.*" -type f -delete

install:
	goimports -w cmd
	go install gitlab.com/cugu/kaitai.go/cmd/...

code:
	# compiler parser/...
	compiler `find . -name "*.ksy" -type f | grep -v "/kaitai/"`

test:
	# go test --bench=. -v gitlab.com/cugu/kaitai.go/parser/...
	go test -v `go list gitlab.com/cugu/kaitai.go/parser/... 2> /dev/null | grep -v "/kaitai/"`
