FROM golang:alpine AS build

WORKDIR /src
COPY . .

RUN apk add --no-cache vips-dev gcc musl-dev
RUN go build

FROM alpine:latest AS bin

RUN apk add --no-cache vips

WORKDIR /app
RUN mkdir /var/cache/librarian
COPY --from=build /src/templates/static /static
COPY --from=build /src/config.example.yml config.yml
COPY --from=build /src/librarian .

EXPOSE 3000

CMD ["/app/librarian"]
