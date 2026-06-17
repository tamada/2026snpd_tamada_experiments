# Experiment 3: Resilience evaluation

This experiment evaluates the similarities in the birthmarks produced by the same software compiled for different architectures and platforms.
We employ \bftt{bzip2} and its go and rust implementations as the target software, which is the same software as in Section~\ref{sect:target_exp2}.
We compile \bftt{bzip2} for macOS Mach-O, Linux ELF, and Windows PE in \texttt{arm64} and \texttt{x86\_64} architectures.
Compilation is performed in GitHub Actions, some Docker images, and a macOS M4 machine (Tahoe 26.3.1).

We use `clang` and `gcc` for C, `go` and [`TinyGo`](https://tinygo.org) for Go, and `rustc` for Rust.
Note that `**bzip2rs**` has two features for compression: a pure Rust implementation and a delegate to `libbz2`.
We should choose one at compile time.
In the Windows environment, we also compile `**bzip2**` with `msvc` (Microsoft Visual C++).
However, due to compilation difficulty in some environments, we were unable to prepare all the binaries.

Besides, the main characteristic of the compilers used is that `clang`, `tinygo`, and `rustc` use LLVM,
while `gcc`, `msvc`, and `go` rely on proprietary technologies.
Furthermore, because `go`'s runtime is very feature-rich, the resulting software tends to be large.

The resultant 35 executables are shown in Table~\ref{table:target_exp3}, and their file sizes are also shown.
All executables are linked statically to the `bzip2` library, and
thus the birthmarks are expected to be similar.
Then, we apply the proposed birthmarking techniques to calculate their similarities.

## Directory layouts

- `birthmarks`: extracted birthmarks of the binaries located on `executables` directory.
- `executables`: the binaries of bzip2, bzip2go, bzip2rs, factorization, md5, and sha256 implemented in C, Go, and Rust.
- `sources`: the source code of factorization, md5, and sha256 implemented in C, Go, and Rust.
- `hungarian`: the comparison results with Hungarian algorithm for aggregation of birthmarks.
  - `images`: the heatmap image from the tables directory data.
  - `results`: the raw comparison results in CSV format, listed the comparison pairs and their similarity scores.
  - `tables`: the comparison results in tabular format.
- `top-n-all`: the comparison results with top-n (n=all) algorithm for aggregation of birthmarks.
  - `images`: the heatmap image from the tables directory data.
  - `results`: the raw comparison results in CSV format, listed the comparison pairs and their similarity scores.
  - `tables`: the comparison results in tabular format.
- `top-n-one`: the comparison results with top-n (n=1) algorithm for aggregation of birthmarks.
  - `images`: the heatmap image from the tables directory data.
  - `results`: the raw comparison results in CSV format, listed the comparison pairs and their similarity scores.
  - `tables`: the comparison results in tabular format.
- `architecture.xlsx`: the data for comparing the similarities across architectures (Section V-D 2 (Result I: Cross architecture), Figure 5).
- `whole_results.xlsx`: The raw comparison results of all pairs, and the extracting data for some analyses.
  - This excel file is related to Section V-D 3 (Result II: clang for Windows), and V-D 4 (Result III: Cross languages, Figure 6)
- `scatter.xlsx`: data for plotting the charts shown in Figure 7 (Section V-D 5, (Result IV: Comparing Algorithms)).
