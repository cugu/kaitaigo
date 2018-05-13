rm -rf src test_formats

go run kaitai.go -d test_formats formats/*.ksy
goimports -w test_formats

mkdir -p src
cp -r test_formats src
cp -r runtime src
mv src/runtime src/kgruntime

# run tests
export GOPATH=$GOPATH:$PWD

ABS_TEST_OUT_DIR="src"
ABS_REPORT_LOG="$ABS_TEST_OUT_DIR/report.log"
keep_compiling=1
while [ "$keep_compiling" = 1 ]; do
    if go test -v test/* >"$ABS_REPORT_LOG" 2>&1; then
        keep_compiling=0
        cat "$ABS_TEST_OUT_DIR/report.log"
    else
        echo "Got error:"
        cat "$ABS_REPORT_LOG"
        if egrep "^src/.*:[0-9][0-9]*:" "$ABS_REPORT_LOG" >"$ABS_TEST_OUT_DIR/err.now"; then
            cat "$ABS_TEST_OUT_DIR/err.now" >>"$ABS_TEST_OUT_DIR/build.fails"
            sed 's/:.*//' <"$ABS_TEST_OUT_DIR/err.now" | sort -u >"$ABS_TEST_OUT_DIR/to_delete.now"
            xargs rm <"$ABS_TEST_OUT_DIR/to_delete.now"
            echo "Trying to recover..."
            keep_compiling=1
        elif grep -q '^=== RUN' "$ABS_REPORT_LOG"; then
            echo "Tests completed partially..."
            keep_compiling=0
        else
            echo "Unable to recover, bailing out :("
            keep_compiling=0
            exit 1
        fi
    fi
done