#!/usr/bin/bash

rm -Rf client
git submodule init
git submodule update
cd client
npm install
npm run build-prod
cd ..

rm -Rf build
mkdir build

for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
        export GOOS
        export GOARCH
        if [ $GOOS = "windows" ]; then
            packr build -v -o build/shortener-$GOOS-$GOARCH.exe
        else
            packr build -v -o build/shortener-$GOOS-$GOARCH
        fi
    done
done

packr clean
