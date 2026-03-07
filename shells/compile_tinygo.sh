#! /bin/sh

function compile_go() {
    GOARCH=$1
    GOOS=$2
    echo "compile bzip2go with tinygo on platform $GOOS/$GOARCH"
    PLATFORM=$GOARCH-$GOOS
    echo "compile go on platform $GOOS/$GOARCH"
    mkdir -p ../build/$PLATFORM/tinygo_O{0,1,2,s,z}
    ext=""
    if [[ "$GOOS" == "windows" ]]; then
        ext=".exe"
    fi

    GOOS=$GOOS GOARCH=$GOARCH tinygo build -o ../build/$PLATFORM/tinygo_O0/bzip2go${ext} -opt=0 ./...
    GOOS=$GOOS GOARCH=$GOARCH tinygo build -o ../build/$PLATFORM/tinygo_O1/bzip2go${ext} -opt=1 ./...
    GOOS=$GOOS GOARCH=$GOARCH tinygo build -o ../build/$PLATFORM/tinygo_O2/bzip2go${ext} -opt=2 ./...
    GOOS=$GOOS GOARCH=$GOARCH tinygo build -o ../build/$PLATFORM/tinygo_Os/bzip2go${ext} -opt=s ./...
    GOOS=$GOOS GOARCH=$GOARCH tinygo build -o ../build/$PLATFORM/tinygo_Oz/bzip2go${ext} -opt=z ./...
}

cd bzip2go
compile_go amd64 linux
compile_go arm64 linux
compile_go amd64 darwin
compile_go arm64 darwin
compile_go amd64 windows
compile_go arm64 windows
