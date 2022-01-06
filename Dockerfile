FROM golang:alpine AS build

WORKDIR /src
COPY . .

ARG GOARCH=amd64
RUN env GOARCH=${GOARCH} go build

FROM alpine:latest as bin

WORKDIR /app
RUN mkdir /var/cache/librarian
COPY --from=build /src/librarian .

EXPOSE 3000

CMD ["sh", "-c", "/app/librarian"]
