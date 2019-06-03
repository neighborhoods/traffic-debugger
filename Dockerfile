FROM golang:1.12-alpine3.9 AS source

RUN apk add --no-cache git

RUN mkdir /src
WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/

FROM source AS build
RUN go build -o /server ./server/

FROM alpine:3.9
COPY --from=build /server /server
RUN chmod +x /server

ENTRYPOINT ["/server"]
