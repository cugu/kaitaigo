all: clean install code test

clean:
	@printf '\nClean\n'
	# find parser -name "*.ksy.*" -type f -delete

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
	@# compiler `find parser/kaitai/type_int_unary_op -name "*.ksy" -type f | grep -v "/enum_fancy/"`
	@# compiler `find parser/kaitai/user_type -name "*.ksy" -type f | grep -v "/enum_fancy/"`
	@# compiler `find parser/kaitai/position_to_end -name "*.ksy" -type f | grep -v "/enum_fancy/"`

test: successful_tests no_tests failing_tests build_failing_tests deprecated_tests
	@# gotestsum --no-summary errors,failed gitlab.com/cugu/kaitai.go/parser/... 2>/dev/null

successful_tests:
	@printf '\n\nTest\n'
	@go test gitlab.com/cugu/kaitai.go/cmd/compiler
	@go test gitlab.com/cugu/kaitai.go/parser/mbr
	@go test gitlab.com/cugu/kaitai.go/parser/apfs
	@go test gitlab.com/cugu/kaitai.go/parser/gpt
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/docstrings
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/enum_0
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
	@go test gitlab.com/cugu/kaitai.go/parser/kaitai/position_to_end # fixed io


no_tests:
	@# go test -v bits_byte_aligned & true
	@# go test -v bits_enum & true
	@# go test -v bits_simple & true
	@# go test -v cast_nested & true
	@# go test -v cast_to_imported & true
	@# go test -v cast_to_top & true
	@# go test -v debug_0 & true
	@# go test -v debug_enum_name & true
	@# go test -v default_endian_expr_exception & true
	@# go test -v default_endian_expr_inherited & true
	@# go test -v enum_1 & true
	@# go test -v enum_deep & true
	@# go test -v enum_deep_literals & true
	@# go test -v enum_for_unknown_id & true
	@# go test -v enum_if & true
	@# go test -v enum_negative & true
	@# go test -v enum_of_value_inst & true
	@# go test -v enum_to_i & true
	@# go test -v eof_exception_bytes & true
	@# go test -v eof_exception_u4 & true
	@# go test -v expr_enum & true
	@# go test -v expr_io_eof & true
	@# go test -v for_rel_imports & true
	@# go test -v if_instances & true
	@# go test -v imports0 & true
	@# go test -v imports_abs & true
	@# go test -v imports_abs_abs & true
	@# go test -v imports_abs_rel & true
	@# go test -v imports_circular_a & true
	@# go test -v imports_circular_b & true
	@# go test -v imports_rel_1 & true
	@# go test -v index_sizes & true
	@# go test -v index_to_param_eos & true
	@# go test -v index_to_param_expr & true
	@# go test -v index_to_param_until & true
	@# go test -v instance_io_user & true
	@# go test -v instance_user_array & true
	@# go test -v ks_path & true
	@# go test -v nav_parent2 & true
	@# go test -v nav_parent3 & true
	@# go test -v nav_parent_false & true
	@# go test -v nav_parent_override & true
	@# go test -v nav_parent_switch & true
	@# go test -v non_standard & true
	@# go test -v opaque_external_type & true
	@# go test -v opaque_external_type_02_child & true
	@# go test -v opaque_external_type_02_parent & true
	@# go test -v opaque_with_param & true
	@# go test -v optional_id & true
	@# go test -v params_call_short & true
	@# go test -v params_def & true
	@# go test -v params_pass_struct & true
	@# go test -v params_pass_usertype & true
	@# go test -v process_coerce_switch & true
	@# go test -v recursive_one & true
	@# go test -v repeat_until_sized & true
	@# go test -v str_literals & true
	@# go test -v switch_bytearray & true
	@# go test -v switch_cast & true
	@# go test -v switch_integers & true
	@# go test -v switch_integers2 & true
	@# go test -v switch_manual_enum & true
	@# go test -v switch_manual_int & true
	@# go test -v switch_manual_int_else & true
	@# go test -v switch_manual_int_size & true
	@# go test -v switch_manual_int_size_else & true
	@# go test -v switch_manual_int_size_eos & true
	@# go test -v switch_manual_str & true
	@# go test -v switch_manual_str_else & true
	@# go test -v switch_multi_bool_ops & true
	@# go test -v switch_repeat_expr & true
	@# go test -v ts_packet_header & true
	@# go test -v type_ternary_opaque & true
	@# go test -v yaml_ints & true

failing_tests:
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

build_failing_tests:
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/bcd_user_type_be & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/bcd_user_type_le & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/buffered_struct & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/docstrings_docref & true # [build failed]
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
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_coerce_bytes & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_coerce_usertype1 & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/process_coerce_usertype2 & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_eos_struct & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/repeat_until_complex & true # [build failed]
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/type_ternary & true # [build failed]

deprecated_tests:
	@# Can be fixed
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/nested_types3 & true # accessing nested types is not allowed
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/enum_fancy & true # no fancy enums

	@# Will not be fixed
	@# go test -v gitlab.com/cugu/kaitai.go/parser/kaitai/nested_same_name2 # dublicate names are not allowed


perf:
	go test -run none -cpuprofile apfsprof -bench . -v gitlab.com/cugu/kaitai.go/parser/apfs
	go tool pprof apfsprof
	# go test -run=none -cpuprofile=mbrprof --bench=. -v gitlab.com/cugu/kaitai.go/parser/mbr
	# go tool pprof mbrprof
	# go test -run=none -cpuprofile=gptprof --bench=. -v gitlab.com/cugu/kaitai.go/parser/gpt
	# go tool pprof gptprof

memory:
	 go test -run none -memprofile apfs.mem.prof  -bench . -v gitlab.com/cugu/kaitai.go/parser/apfs