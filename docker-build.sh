#!/bin/bash

version=$1

tag="korylprince/url-shortener-server:$version"

docker build --no-cache --build-arg "VERSION=$version" --tag "$tag" .

docker push "$tag"
