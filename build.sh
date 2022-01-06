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

docker build -t nineteengladespool/librarian:amd64 --build-arg GOARCH=amd64 .
docker push nineteengladespool/librarian:amd64
docker build -t nineteengladespool/librarian:arm64 --build-arg GOARCH=arm64 .
docker push nineteengladespool/librarian:arm64
docker manifest create nineteengladespool/librarian:latest --amend nineteengladespool/librarian:amd64 --amend nineteengladespool/librarian:arm64 
docker manifest create nineteengladespool/librarian:${VERSION} --amend nineteengladespool/librarian:amd64 --amend nineteengladespool/librarian:arm64
docker manifest push nineteengladespool/librarian:latest
docker manifest push nineteengladespool/librarian:${VERSION}