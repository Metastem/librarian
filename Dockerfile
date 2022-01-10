FROM golang:alpine AS build

WORKDIR /src
RUN apk --no-cache add git
RUN git clone https://codeberg.org/librarian/librarian .

RUN go build

FROM alpine:latest as bin

WORKDIR /app
COPY --from=build /src/librarian .

EXPOSE 3000

CMD ["/app/librarian"]
