#!/usr/bin/bash

rm -Rf build
mkdir build

for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
        export GOOS
        export GOARCH
        if [ $GOOS = "windows" ]; then
            go build -v -o build/shortener-$GOOS-$GOARCH.exe
        else
            go build -v -o build/shortener-$GOOS-$GOARCH
        fi
    done
done
