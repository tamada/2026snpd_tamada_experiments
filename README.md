# 2026snpd_tamada

## Target programs and compiler options

- [x] [bzip2](https://sourceware.org/bzip2/)
  - `gcc`: `O0`, `O1`, `O2`, `O3`, `Os`, `Oz`
  - `clang`: `O0`, `O1`, `O2`, `O3`, `Os`, `Oz`
- [x] [pedroalbanese/bzip2](https://github.com/pedroalbanese/bzip2), written in Go.
  - Go compiler: `default`, `-ldflags="-s -w"`
  - TinyGo: `-opt=0`, `-opt=1`, `-opt=2`, `-opt=s`, `-opt=z`
- [x] [tamada/bzip2rs](https://github.com/tamada/bzip2rs), written in Rust (with pure Rust crate, and libbzip2 wrapper crate).
  - `default` feature (pure Rust crate): `debug`, `release`
    - `opt-level`: `0`, `1`, `2`, `3`, `s`, `z` (`RUSTFLAGS="-C opt-level=0" cargo build`, etc.)
  - `sys` feature (libbzip2 wrapper crate): `debug`, `release`
    - `opt-level`: `0`, `1`, `2`, `3`, `s`, `z` (`RUSTFLAGS="-C opt-level=0" cargo build --features=sys`, etc.)

6 optimization levels, 5 compilers (3 C/C++ compilers, 1 Rust compiler, 1 Zig compiler), resulting in 30 configurations per target program.

## Format

- ELF (Linux)
  - `x86_64`
  - `aarch64`
- MachO (macOS)
  - `x86_64`
  - `aarch64`
- PE (Windows)
  - `x86_64`
  - `aarch64`

6 platforms, resulting in 180 configurations per target program.

## How to compile

### :clown_face: Go

Go is quite easy to cross-compile by setting `GOOS` and `GOARCH` environment variables.
Therefore, to compile `go-bzip2` to all 6 platforms, we can run `shells/compile_go.sh` and `shells/compile_tinygo.sh`.

### :musl: Rust

In the macOS environment, we install the Rust targets for x86-64-apple-darwin, aarch64-apple-darwin, aarch64-pc-windows-msvc, and x86-64-pc-windows-msvc targets, cargo component xwin, as follows:

```sh
rustup target add rust-std-x86_64-apple-darwin     \
                  rust-std-aarch64-apple-darwin    \
                  rust-std-aarch64-pc-windows-msvc \
                  rust-std-x86_64-pc-windows-msvc
cargo install cargo-xwin
```

Then, it uses Docker environment to compile to arm64-linux and amd64-linux platform, it uses `compile.sh` and `shells/compile_rust.sh`.

### :smile: Clang

Clang is cross-compile friendly, and we can specify the target platform by `--target` option.

### :face_screaming_in_fear: Gcc

