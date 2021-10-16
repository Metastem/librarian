FROM golang:alpine AS build

WORKDIR /src
COPY . .

ARG opts
RUN env ${opts} go build

FROM alpine:latest as bin

WORKDIR /app
RUN mkdir /var/cache/librarian
COPY --from=build /src/templates/static /static
COPY --from=build /src/config.example.yml config.yml
COPY --from=build /src/librarian .

EXPOSE 3000

CMD ["/app/librarian"]
