#!/bin/bash

version=$1

tag="korylprince/url-shortener-server:$version"

docker build --no-cache --build-arg "VERSION=$version" --tag "$tag" .

docker push "$tag"

if [ "$2" = "latest" ]; then
    docker tag "$tag" "korylprince/url-shortener-server:latest"
    docker push "korylprince/url-shortener-server:latest"
fi
