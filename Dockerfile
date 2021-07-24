FROM golang:alpine AS build

WORKDIR /src
COPY . .

RUN go build

FROM alpine:latest AS bin

WORKDIR /app
COPY --from=build /src/config.example.yml config.yml
COPY --from=build /src/librarian .

EXPOSE 3000

CMD ["/app/librarian"]