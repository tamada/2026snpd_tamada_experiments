#! /bin/sh

function compile_go() {
    echo "compile bzip2go with tinygo on platform $GOOS/$GOARCH"
    GOARCH=$1
    GOOS=$2
    PLATFORM=$GOARCH-$GOOS
    echo "compile go on platform $GOOS/$GOARCH"
    mkdir -p build/$PLATFORM/tinygo_O{0,1,2,s,z}

    GOOS=$GOOS GOARCH=$GOARCH tinygo build -o ../build/$PLATFORM/tinygo_O0/bzip2go -opt=0 ./...
    GOOS=$GOOS GOARCH=$GOARCH tinygo build -o ../build/$PLATFORM/tinygo_O1/bzip2go -opt=1 ./...
    GOOS=$GOOS GOARCH=$GOARCH tinygo build -o ../build/$PLATFORM/tinygo_O2/bzip2go -opt=2 ./...
    GOOS=$GOOS GOARCH=$GOARCH tinygo build -o ../build/$PLATFORM/tinygo_Os/bzip2go -opt=s ./...
    GOOS=$GOOS GOARCH=$GOARCH tinygo build -o ../build/$PLATFORM/tinygo_Oz/bzip2go -opt=z ./...
}

cd go-bzip2
compile_go amd64 linux
compile_go arm64 linux
compile_go amd64 darwin
compile_go arm64 darwin
compile_go amd64 windows
compile_go arm64 windows
