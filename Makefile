all: clean install code test

clean:
	@printf '\nClean\n'
	find . -name "*.ksy.*" -type f -delete

install:
	@printf '\n\nInstall\n'
	goimports -w cmd/compiler
	goimports -w parser
	goimports -w cmd
	goimports -w runtime
	go install gitlab.com/cugu/kaitai.go/cmd/...

code:
	@printf '\n\nCode\n'
	@# compiler parser/...
	@#compiler `find . -name "*.ksy" -type f` # | grep -v "/kaitai/"`
	compiler `find "cmd/compiler" -name "*.ksy" -type f`
	compiler `find "parser/mbr" -name "*.ksy" -type f`
	compiler `find "parser/gpt" -name "*.ksy" -type f`
	compiler `find "parser/kaitai/hello_world" -name "*.ksy" -type f`
	compiler `find "parser/kaitai/expr_0" -name "*.ksy" -type f`

test:
	@printf '\n\nTest\n'
	@# go test  -v gitlab.com/cugu/kaitai.go/parser/...
	@go test gitlab.com/cugu/kaitai.go/cmd/compiler
	@go test gitlab.com/cugu/kaitai.go/parser/mbr
	@go test gitlab.com/cugu/kaitai.go/parser/gpt
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/hello_world
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/expr_0


perf:
	go test -run none -cpuprofile apfsprof -bench . -v gitlab.com/cugu/kaitai.go/parser/apfs
	go tool pprof apfsprof
	# go test -run=none -cpuprofile=mbrprof --bench=. -v gitlab.com/cugu/kaitai.go/parser/mbr
	# go tool pprof mbrprof
	# go test -run=none -cpuprofile=gptprof --bench=. -v gitlab.com/cugu/kaitai.go/parser/gpt
	# go tool pprof gptprof

memory:
	 go test -run none -memprofile apfs.mem.prof  -bench . -v gitlab.com/cugu/kaitai.go/parser/apfs