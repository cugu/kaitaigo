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

	@# build fail
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/bcd_user_type_be & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/bcd_user_type_le & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/buffered_struct & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/docstrings_docref & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/enum_0 & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/enum_fancy & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/expr_2 & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/expr_3 & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/expr_array & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/expr_bytes_cmp & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/expr_io_pos & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/fixed_contents & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/fixed_struct & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/float_to_i & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/floating_points & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/if_values & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/integers & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/nav_parent_vs_value_inst & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/nested_same_name2 & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/nested_types3 & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/position_to_end & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_coerce_bytes & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_coerce_usertype1 & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_coerce_usertype2 & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_eos_struct & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_until_complex & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/type_ternary & true # [build failed]

	@# fail
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/default_endian_mod & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/bytes_pad_term  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/default_big_endian  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/expr_mod  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/instance_std_array  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/if_struct  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/position_in_seq  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/position_abs  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_rotate  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_custom  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_to_user  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_xor4_const  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_xor4_value  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_xor_const  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_xor_value  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_eos_u4  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_n_strz  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_n_strz_double  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/str_encodings  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_until_s4  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/str_encodings_default  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/str_pad_term  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/term_bytes  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/str_pad_term_empty  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/term_strz  & true
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/zlib_with_header_78  & true

	@go test gitlab.com/cugu/kaitai.go/parser/mbr
	@go test gitlab.com/cugu/kaitai.go/parser/apfs
	@go test gitlab.com/cugu/kaitai.go/parser/gpt
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/docstrings
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/expr_0
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/expr_1
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/hello_world
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/instance_std
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/js_signed_right_shift
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/meta_xref
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/multiple_use
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/nav_parent
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/nav_parent_false2
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/nav_root
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/nested_same_name
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/nested_types
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/nested_types2
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_n_struct
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/str_eos
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/str_literals2
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/type_int_unary_op
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/user_type

	@# gotestsum --no-summary errors,failed gitlab.com/cugu/kaitai.go/parser/... 2>/dev/null




perf:
	go test -run none -cpuprofile apfsprof -bench . -v gitlab.com/cugu/kaitai.go/parser/apfs
	go tool pprof apfsprof
	# go test -run=none -cpuprofile=mbrprof --bench=. -v gitlab.com/cugu/kaitai.go/parser/mbr
	# go tool pprof mbrprof
	# go test -run=none -cpuprofile=gptprof --bench=. -v gitlab.com/cugu/kaitai.go/parser/gpt
	# go tool pprof gptprof

memory:
	 go test -run none -memprofile apfs.mem.prof  -bench . -v gitlab.com/cugu/kaitai.go/parser/apfs