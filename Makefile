all: clean install code test

clean:
	@printf '\nClean\n'
	find . -name "*.ksy.*" -type f -delete

install:
	@printf '\n\nInstall\n'
	goimports -w cmd
	go install gitlab.com/cugu/kaitai.go/cmd/...

code:
	@printf '\n\nCode\n'
	@# compiler parser/...
	compiler `find . -name "*.ksy" -type f | grep -v "/kaitai/"`

test:
	@printf '\n\nTest\n'
	@# go test --bench=. -v gitlab.com/cugu/kaitai.go/parser/...
	go test gitlab.com/cugu/kaitai.go/cmd/compiler
	go test -v `go list gitlab.com/cugu/kaitai.go/parser/... 2> /dev/null | grep -v "/kaitai/"`
