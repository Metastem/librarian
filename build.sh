#!/bin/sh
GOOS=linux GOARCH=amd64 go build
tar -cf librarian-${VERSION}-linux-amd64.tar.gz librarian
rm librarian

GOOS=linux GOARCH=arm64 go build
tar -cf librarian-${VERSION}-linux-arm64.tar.gz librarian
rm librarian

GOOS=openbsd GOARCH=amd64 go build
tar -cf librarian-${VERSION}-openbsd-amd64.tar.gz librarian
rm librarian

sudo docker buildx create --use
sudo docker buildx build --platform linux/amd64,linux/arm64 -t nineteengladespool/librarian:latest -t nineteengladespool/librarian:${VERSION} --push .