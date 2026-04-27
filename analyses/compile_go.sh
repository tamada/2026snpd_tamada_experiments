#! /bin/sh

function compile_go() {
    GOARCH=$1
    GOOS=$2
    echo "compile bzip2go with Go on platform $GOOS/$GOARCH"
    PLATFORM=$GOARCH-$GOOS
    mkdir -p ../build/$PLATFORM/go_{default,sw}
    ext=""
    if [[ "$GOOS" == "windows" ]]; then
        ext=".exe"
    fi

    GOOS=$GOOS GOARCH=$GOARCH go build -o ../build/$PLATFORM/go_default/bzip2go${ext} ./...
    GOOS=$GOOS GOARCH=$GOARCH go build -o ../build/$PLATFORM/go_sw/bzip2go${ext} -ldflags="-s -w" ./...
}

cd bzip2go
compile_go amd64 linux
compile_go arm64 linux
compile_go amd64 darwin
compile_go arm64 darwin
compile_go amd64 windows
compile_go arm64 windows
