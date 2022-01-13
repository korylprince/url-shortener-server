#!/usr/bin/bash

TAG=$1

TEMP=$(mktemp -d)
trap 'rm -rf -- "$TEMP"' EXIT

pushd $TEMP
git clone https://github.com/korylprince/url-shortener-client.git .
git checkout $TAG
NODE_VERSION=lts/fermium $HOME/.nvm/nvm-exec npm install
NODE_VERSION=lts/fermium $HOME/.nvm/nvm-exec npm run build-prod

popd
rm -rf client
cp -r $TEMP/dist client
