#!/bin/bash -e

VERSION="0.1.0-test2"

rm -rf dockerbuild || true
mkdir dockerbuild

cp Dockerfile dockerbuild/Dockerfile-amd64
cp Dockerfile dockerbuild/Dockerfile-arm
cp Dockerfile dockerbuild/Dockerfile-arm64

sed -E 's|FROM alpine|FROM amd64/alpine|' -i dockerbuild/Dockerfile-amd64
sed -E 's/GOARCH=/GOARCH=amd64/' -i dockerbuild/Dockerfile-amd64

sed -E 's|FROM alpine|FROM arm32v7/alpine|'   -i dockerbuild/Dockerfile-arm
sed -E 's/GOARCH=/GOARCH=arm/'   -i dockerbuild/Dockerfile-arm
# sed -E 's/CC "\$CC"/CC "arm-linux-gnueabi-gcc"/'  -i dockerbuild/Dockerfile-arm
# sed -E 's/build-essential/gcc-arm* build-essential/'  -i dockerbuild/Dockerfile-arm

sed -E 's|FROM alpine|FROM arm64v8/alpine|' -i dockerbuild/Dockerfile-arm64
sed -E 's/GOARCH=/GOARCH=arm64/' -i dockerbuild/Dockerfile-arm64
#aarch64-linux-musl-gcc

cat dockerbuild/Dockerfile-arm

docker build -f dockerbuild/Dockerfile-amd64 -t sequentialread/gotosocial:$VERSION-amd64 .
docker build -f dockerbuild/Dockerfile-arm   -t sequentialread/gotosocial:$VERSION-arm .
docker build -f dockerbuild/Dockerfile-arm64 -t sequentialread/gotosocial:$VERSION-arm64 .

docker push sequentialread/gotosocial:$VERSION-amd64
docker push sequentialread/gotosocial:$VERSION-arm
docker push sequentialread/gotosocial:$VERSION-arm64

export DOCKER_CLI_EXPERIMENTAL=enabled


docker manifest create  sequentialread/gotosocial:$VERSION \
   sequentialread/gotosocial:$VERSION-amd64 \
   sequentialread/gotosocial:$VERSION-arm \
   sequentialread/gotosocial:$VERSION-arm64 

docker manifest annotate --arch amd64 sequentialread/gotosocial:$VERSION sequentialread/gotosocial:$VERSION-amd64
docker manifest annotate --arch arm sequentialread/gotosocial:$VERSION sequentialread/gotosocial:$VERSION-arm
docker manifest annotate --arch arm64 sequentialread/gotosocial:$VERSION sequentialread/gotosocial:$VERSION-arm64

docker manifest push sequentialread/gotosocial:$VERSION

rm -rf dockerbuild || true