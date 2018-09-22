**ksy**
replace "package *" with "package spec"

**tests**

remove "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
remove . "test_formats"
replace "../../src/" with "../../../testdata/kaitai"
remove "s := kaitai.NewStream(f)"
replace "err = r.Read(s, &r, &r)" with "err = r.Decode(f)"