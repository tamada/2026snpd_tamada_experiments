# Experiment 3: Resilience evaluation

This experiment evaluates the similarities in the birthmarks produced by the same software compiled for different architectures and platforms.
We employ `**bzip2**` and its Go and Rust implementations as the target software, which is the same software as in Section V-C.
We compile `**bzip2**` for macOS Mach-O, Linux ELF, and Windows PE in amd64 and arm64 architectures.
Note that the number of functions in the source code is: 129 (`**bzip2**`), 5 (`**bzip2go**`), and 22 (`**bzip2rs**`) via `ctags`.

We prepare the executables using `clang` and `gcc` for C, `go` and [`TinyGo`](https://tinygo.org) for Go, and `rustc` for Rust.
Note that `bzip2rs` has two features for compression: a pure Rust implementation and a delegate to `libbz2`.
We chose one at compile time.
In the Windows environment, we also use `msvc` (Microsoft Visual C++) for \bftt{bzip2}.
%
Almost all compilations were performed in [GitHub Actions](https://github.com/tamada/2026snpd_tamada_experiments).
Also, we use Docker images (golang:1.26.3-bookworm and tinygo/tinygo:0.40.0) to compile \bftt{bzip2go} across all three platforms.
Compiling \bftt{bzip2rs} was performed on the local macOS machine with [zigbuild](https://github.com/rust-cross/cargo-zigbuild) for linking.
Specifically, for the Windows platform, we use [xwin](https://github.com/rust-cross/cargo-xwin) for `**bzip2rs**` on macOS,
with target `pc-windows-msvc` for amd64 and arm64 architectures.
Unfortunately, we could not compile `**bzip2go**` with `TinyGo` on Windows platform because it does not support Windows well.
In addition, in macOS environment, `gcc` is actually `clang`.
Therefore, we should install it explicitly to use the actual `gcc` (see the version information of \texttt{gcc} on macOS environment.).

The resulting 36 executables are shown in Table III, which
also includes their file sizes, the compilers used, their versions, and the number of functions in the Oinkie-IR (birthmarks).
The number of `**bzip2**` functions is reduced after compiling and lifting, which may cause inline expansion or unused functions.

Besides, `clang`, `TinyGo`, and `rustc` use LLVM,
while `gcc`, `msvc`, and `go` rely on distinct backend technologies.
This is reflected in the file sizes in Table III.
Go's feature-rich runtime yields the largest executables, followed by TinyGo's garbage-collected runtime.
Rust's thin runtime and C's low-level nature result in smaller binaries.

Ideally, the compiler versions should be the same across all platforms,
however, due to the complexity of preparing, we could not achieve that.
Then, we apply the proposed birthmarking techniques to calculate their similarities.

## Directory layouts

- `birthmarks`: extracted birthmarks of the binaries located on `executables` directory.
- `executables`: the binaries of bzip2, bzip2go, bzip2rs, factorization, md5, and sha256 implemented in C, Go, and Rust.
- `sources`: the source code of factorization, md5, and sha256 implemented in C, Go, and Rust.
- `hungarian`: the comparison results with Hungarian algorithm for aggregation of birthmarks.
  - `images`: the heatmap image from the tables directory data.
  - `results`: the raw comparison results in CSV format, listed the comparison pairs and their similarity scores.
  - `tables`: the comparison results in tabular format.
