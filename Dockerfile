FROM golang:alpine AS build

WORKDIR /src
COPY . .

ARG GOARCH=amd64
RUN env GOARCH=${GOARCH} go build
RUN mkdir /var/cache/librarian

FROM scratch as bin

WORKDIR /app
COPY --from=build /var/cache/librarian /var/cache/librarian
COPY --from=build /src/librarian .

EXPOSE 3000

CMD ["/app/librarian"]
