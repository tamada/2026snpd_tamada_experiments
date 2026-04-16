#! /bin/bash

echo "compile bzip2go with Go"
docker run -it --rm -v $PWD:/opt -w /opt golang:bookworm      bash scripts/compile_go.sh

echo "compile go-bzip2 with TinyGo"
docker run -it --rm -v $PWD:/opt -w /opt tinygo/tinygo:latest bash scripts/compile_tinygo.sh

echo "compile bzip2 written in the C with clang in arm64/linux"
docker run -it --rm -v $PWD:/opt -w /opt silkeh/clang:latest  bash scripts/compile_with_clang.sh arm64-linux

echo "compile bzip2 written in the C with gcc in arm64/linux"
docker run -it --rm -v $PWD:/opt -w /opt gcc:13-bookworm      bash scripts/compile_with_gcc.sh arm64-linux

echo "compile bzip2rs written in the Rust in arm64/linux"
docker run -it --rm -v $PWD:/opt -w /opt rust:bookworm        bash scripts/compile_rust.sh arm64-linux

echo "compile bzip2 written in the C with clang in arm64/linux"
docker run -it --rm -v $PWD:/opt -w /opt --platform linux/amd64 silkeh/clang:latest bash scripts/compile_with_clang.sh amd64-linux

echo "compile bzip2 written in the C with gcc in arm64/linux"
docker run -it --rm -v $PWD:/opt -w /opt --platform linux/amd64 gcc:13-bookworm     bash scripts/compile_with_gcc.sh   amd64-linux

echo "compile bzip2rs written in the Rust in amd64/linux"
docker run -it --rm -v $PWD:/opt -w /opt --platform linux/amd64 rust:bookworm       bash scripts/compile_rust.sh       amd64-linux
