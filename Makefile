all: clean install code test

clean:
	@printf '\nClean\n'
	find parser -name "*.ksy.*" -type f -delete

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
	@compiler `find parser -name "*.ksy" -type f | grep -v "/enum_fancy/"`

test:
	@printf '\n\nTest\n'
	@go test gitlab.com/cugu/kaitai.go/cmd/compiler
	@go test gitlab.com/cugu/kaitai.go/parser/...
	@# go test gitlab.com/cugu/kaitai.go/parser/mbr
	@# go test gitlab.com/cugu/kaitai.go/parser/gpt
	@# go test gitlab.com/cugu/kaitai.go/parser/apfs
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/expr_0
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/expr_1
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/hello_world
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/nav_parent
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/nested_types
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/nested_types2
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/str_eos

	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/buffered_struct
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/expr_2
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/nested_types3
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_eos_struct
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_n_strz
	@# go test gitlab.com/cugu/kaitai.go/parser/kaitai/str_encodings


perf:
	go test -run none -cpuprofile apfsprof -bench . -v gitlab.com/cugu/kaitai.go/parser/apfs
	go tool pprof apfsprof
	# go test -run=none -cpuprofile=mbrprof --bench=. -v gitlab.com/cugu/kaitai.go/parser/mbr
	# go tool pprof mbrprof
	# go test -run=none -cpuprofile=gptprof --bench=. -v gitlab.com/cugu/kaitai.go/parser/gpt
	# go tool pprof gptprof

memory:
	 go test -run none -memprofile apfs.mem.prof  -bench . -v gitlab.com/cugu/kaitai.go/parser/apfs