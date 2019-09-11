all: clean test install generate_code ks_tests

ex:
	go install dfir.software/kaitaigo/cmd/...
	kaitaigo example/my_format.ksy
	cd example && go build . && ./example

dep:
	@printf '\nGet dependencies\n'
	go mod download

clean:
	@printf '\nClean\n'
	find tests -name "*.ksy.*" -type f -delete
	goimports -w cmd/kaitaigo
	goimports -w runtime

test:
	go test dfir.software/kaitaigo/cmd/...

lint:
	@printf '\nLint\n'
	go get -u gopkg.in/alecthomas/gometalinter.v2
	gometalinter.v2 --install
	gometalinter.v2 cmd/...
	gometalinter.v2 runtime/...

build:
	@printf '\n\nBuild\n'
	go build dfir.software/kaitaigo/cmd/...

install:
	@printf '\n\nInstall\n'
	go install dfir.software/kaitaigo/cmd/...

generate_code:
	@printf '\n\nCode\n'
	kaitaigo `find tests -name "*.ksy" -type f | grep -v "/enum_fancy/"`

ks_tests:
	@printf '\n\nTest\n'
	go test dfir.software/kaitaigo/tests/kaitai/buffered_struct \
		dfir.software/kaitaigo/tests/kaitai/bcd_user_type_be \
		dfir.software/kaitaigo/tests/kaitai/bcd_user_type_le \
		dfir.software/kaitaigo/tests/kaitai/bytes_pad_term \
		dfir.software/kaitaigo/tests/kaitai/default_big_endian \
		dfir.software/kaitaigo/tests/kaitai/docstrings \
		dfir.software/kaitaigo/tests/kaitai/docstrings_docref \
		dfir.software/kaitaigo/tests/kaitai/enum_0 \
		dfir.software/kaitaigo/tests/kaitai/expr_0 \
		dfir.software/kaitaigo/tests/kaitai/expr_1 \
		dfir.software/kaitaigo/tests/kaitai/expr_2 \
		dfir.software/kaitaigo/tests/kaitai/fixed_contents \
		dfir.software/kaitaigo/tests/kaitai/fixed_struct \
		dfir.software/kaitaigo/tests/kaitai/float_to_i \
		dfir.software/kaitaigo/tests/kaitai/hello_world \
		dfir.software/kaitaigo/tests/kaitai/if_struct \
		dfir.software/kaitaigo/tests/kaitai/if_values \
		dfir.software/kaitaigo/tests/kaitai/instance_std \
		dfir.software/kaitaigo/tests/kaitai/instance_std_array \
		dfir.software/kaitaigo/tests/kaitai/integers \
		dfir.software/kaitaigo/tests/kaitai/js_signed_right_shift \
		dfir.software/kaitaigo/tests/kaitai/meta_xref \
		dfir.software/kaitaigo/tests/kaitai/multiple_use \
		dfir.software/kaitaigo/tests/kaitai/nav_parent \
		dfir.software/kaitaigo/tests/kaitai/nav_parent_false2 \
		dfir.software/kaitaigo/tests/kaitai/nav_root \
		dfir.software/kaitaigo/tests/kaitai/nested_same_name \
		dfir.software/kaitaigo/tests/kaitai/nested_types \
		dfir.software/kaitaigo/tests/kaitai/nested_types2 \
		dfir.software/kaitaigo/tests/kaitai/position_abs \
		dfir.software/kaitaigo/tests/kaitai/position_in_seq \
		dfir.software/kaitaigo/tests/kaitai/position_to_end \
		dfir.software/kaitaigo/tests/kaitai/process_coerce_bytes \
		dfir.software/kaitaigo/tests/kaitai/process_custom \
		dfir.software/kaitaigo/tests/kaitai/process_rotate \
		dfir.software/kaitaigo/tests/kaitai/process_xor4_const \
		dfir.software/kaitaigo/tests/kaitai/process_xor4_value \
		dfir.software/kaitaigo/tests/kaitai/process_xor_const \
		dfir.software/kaitaigo/tests/kaitai/process_xor_value \
		dfir.software/kaitaigo/tests/kaitai/repeat_eos_struct \
		dfir.software/kaitaigo/tests/kaitai/repeat_eos_u4 \
		dfir.software/kaitaigo/tests/kaitai/repeat_n_struct \
		dfir.software/kaitaigo/tests/kaitai/repeat_n_strz \
		dfir.software/kaitaigo/tests/kaitai/repeat_n_strz_double \
		dfir.software/kaitaigo/tests/kaitai/repeat_until_complex \
		dfir.software/kaitaigo/tests/kaitai/repeat_until_s4 \
		dfir.software/kaitaigo/tests/kaitai/str_eos \
		dfir.software/kaitaigo/tests/kaitai/str_literals2 \
		dfir.software/kaitaigo/tests/kaitai/str_pad_term \
		dfir.software/kaitaigo/tests/kaitai/str_pad_term_empty \
		dfir.software/kaitaigo/tests/kaitai/term_bytes \
		dfir.software/kaitaigo/tests/kaitai/term_strz \
		dfir.software/kaitaigo/tests/kaitai/type_int_unary_op \
		dfir.software/kaitaigo/tests/kaitai/user_type \
		dfir.software/kaitaigo/tests/kaitai/zlib_with_header_78
	@# Changes
	@# position_to_end # fixed io
	@# repeat_n_strz # fix tests
	@# repeat_n_strz_double # fix tests
	@# repeat_eos_u4 # fix tests, typecast
	@# repeat_until_s4 # fix tests, typecast
	@# repeat_until_complex # fix tests, typecast
	@# if_values # change tests for nil

missing_tests:
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
	@# Could be fixed
	@# go test -v dfir.software/kaitaigo/tests/kaitai/nested_types3 		# accessing nested types is not allowed
	@# go test -v dfir.software/kaitaigo/tests/kaitai/enum_fancy 			# no fancy enums
	@# go test -v dfir.software/kaitaigo/tests/kaitai/default_endian_mod 	# no nested endianess
	@# go test -v dfir.software/kaitaigo/tests/kaitai/str_encodings 		# no other encoding
	@# go test -v dfir.software/kaitaigo/tests/kaitai/str_encodings_default # no other encoding
	@# go test -v dfir.software/kaitaigo/tests/kaitai/expr_bytes_cmp 		# compare []byte
	@# go test -v dfir.software/kaitaigo/tests/kaitai/expr_array 			# need generic min, max funcs
	@# go test -v dfir.software/kaitaigo/tests/kaitai/nav_parent_vs_value_inst # fix type inference

	@# Hard to fix
	@# go test -v dfir.software/kaitaigo/tests/kaitai/expr_mod  			# -2 % 8 => -2
	@# go test -v dfir.software/kaitaigo/tests/kaitai/expr_3 			 	# string compare
	@# go test -v dfir.software/kaitaigo/tests/kaitai/type_ternary 			# xor only on bytes
	@# go test -v dfir.software/kaitaigo/tests/kaitai/process_to_user 		# rol only on bytes
	@# go test -v dfir.software/kaitaigo/tests/kaitai/process_coerce_usertype1 # xor only on bytes
	@# go test -v dfir.software/kaitaigo/tests/kaitai/process_coerce_usertype2 # xor only on bytes
	@# go test -v dfir.software/kaitaigo/tests/kaitai/floating_points 		# float + int does not work
	@# go test -v dfir.software/kaitaigo/tests/kaitai/expr_io_pos 			# size: _io.size - _io.pos ??

	@# Will not be fixed
	@# go test -v dfir.software/kaitaigo/tests/kaitai/nested_same_name2 	# dublicate names are not allowed
