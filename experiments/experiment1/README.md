# Experiment 1: Credibility Evaluation

同じプラットフォームで異なるバイナリが異なることを確認する。

以下の23のバイナリを対象に相互比較を行った。23C2で253の比較を行った。

- Bzip2 (compiled by clang with the default optimization)
- Bzip2 (compiled by gcc with the default optimization)
- Bzip2go (compiled by go build)
- Bzip2go (compiled by tinygo)
- Bzip2rs (compiled by rustc with debug option and use libbz2)
- Bzip2rs (compiled by rustc with release option and use libbz2)
- Bzip2rs (compiled by rustc with debug option and use pure rust bzip2 crate)
- Bzip2rs (compiled by rustc with release option and use pure rust bzip2 crate)
- 素因数分解 (clang, gcc, go, tinygo, rs)
- MD5 (clang, gcc, go, tinygo, rs)
- SHA256 (clang, gcc, go, tinygo, rs)

比較手法は次の8通り。

- Cosine
- Dice
- Euclidean
- Jaccard
- LCS
- Levenshtein
- Simpson
- Weighted Jaccard