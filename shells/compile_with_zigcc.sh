#! /bin/sh

function compile_bzip2_each() {
    CC=$1
    OP=$2
    PLATFORM=$3
    make clean
    make CC=\"$CC\" CFLAGS="-Wall -Winline -$OP -g -D_FILE_OFFSET_BITS=64"
    cp bzip2 ../build/$3/zig-cc_$OP
}

function compile_bzip2() {
    echo "compile bzip2 on platform $1"

    for optimization in O0 O1 O2 O3 Os Oz
    do
        compile_bzip2_each "zig cc" $optimization $1
    done
}

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <platform>"
    exit 1
fi

mkdir -p build/$1/zig-cc_O{0,1,2,3,s,z}

cd bzip2
compile_bzip2 $1
