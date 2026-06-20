# Experiment 2: Credibility Evaluation

> [!WARNING]
> This page is not used, because the compilations were not correct.
> See [`../experiment2-2`](../experiment2-2/)

This experiment evaluates differences of the birthmarks from the different software.
For this purpose, we should choose software with the same purpose and different authors.
Therefore, we first selected [**bzip2**](https://sourceware.org/bzip2) as the target software.
We also chose a Go implementation of bzip2 ([**bzip2go**](https://github.com/pedroalbanese/bzip2)).
In addition, the author implemented a Rust version of bzip2 ([**bzip2rs**](https://github.com/tamada/bzip2rs)) because none existed.

On the other hand, we evaluate differences in birthmarks across programming language levels.
For this, we should use the same algorithm in different languages.
Therefore, we requested for generative AI, Google Gemini to implement factorization, MD5, and SHA256 in C, Go, and Rust.
The prompt to generative AI was "Write the ***SPEC*** logic in ***LANGUAGE*** without standard libraries."
The italic words were replaced with the corresponding specification and language.
Of couse, there are many other famous programming languages.
However, we chose C, Go, and Rust, since they can generate executables and are easy to implement.

Moreover, the algorithms of **bzip2**, **md5**, and **sha256** use many bit calculations,
and thus are expected to have similar instruction sets and high similarities.
On the other hand, the algorithm of **factorization** is implemented by trial division,
which uses multiplication and remainder calculations, and no bit calculations.
Therefore, it is expected to have a different instruction set and low similarities from the other three software.
Next, we compile the obtained four software in three languages to arm64 Mach-O executables on macOS.
We employ `gcc` for C, `go` for Go, and `rustc` for Rust.
No compile options were specified in any cases.
Finally, we obtain 21 executables.
Then, we apply the proposed birthmarking techniques to calculate the similarities among them.

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
- `freq.xlsx`: the frequencies of similarities.
- `ecdf.xlsx`: ECDF analysis of similarities.
