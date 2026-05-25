#! /bin/sh

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <platform>"
    exit 1
fi

builder=cargo
exts=""
if [[ "$1" == "amd64-linux" ]]; then
    target="x86_64-unknown-linux-gnu"
elif [[ "$1" == "arm64-linux" ]]; then
    target="aarch64-unknown-linux-gnu"
elif [[ "$1" == "amd64-darwin" ]]; then
    target="x86_64-apple-darwin"
elif [[ "$1" == "arm64-darwin" ]]; then
    target="aarch64-apple-darwin"
elif [[ "$1" == "amd64-windows" ]]; then
    target="x86_64-pc-windows-msvc"
    builder="cargo xwin"
    exts=".exe"
elif [[ "$1" == "arm64-windows" ]]; then
    target="aarch64-pc-windows-msvc"
    builder="cargo xwin"
    exts=".exe"
else
    echo "Unsupported platform: $1"
    exit 1
fi

function compile_rust_each() {
    echo $PWD, $1, $2
    RUSTFLAGS="-C opt-level=$1" $builder build --release --target $target
    cp target/${target}/release/bzip2rs${exts} ../build/$2/rust_pure_release_O$1 && cargo clean
    RUSTFLAGS="-C opt-level=$1" $builder build --target $target
    cp target/${target}/debug/bzip2rs${exts}   ../build/$2/rust_pure_debug_O$1 && cargo clean

    RUSTFLAGS="-C opt-level=$1" $builder build --release  --features sys --target $target
    cp target/${target}/release/bzip2rs${exts} ../build/$2/rust_sys_release_O$1 && cargo clean
    RUSTFLAGS="-C opt-level=$1" $builder build --features sys --target $target
    cp target/${target}/debug/bzip2rs${exts}   ../build/$2/rust_sys_debug_O$1 && cargo clean
}

function compile_rust() {
    echo "compile rust on platform $1"

    compile_rust_each 0 $1
    compile_rust_each 1 $1
    compile_rust_each 2 $1
    compile_rust_each 3 $1
    compile_rust_each s $1
    compile_rust_each z $1
}

mkdir -p build/$1/rust_{pure,sys}_{debug,release}_O{0,1,2,3,s,z}

cd bzip2rs
compile_rust $1
