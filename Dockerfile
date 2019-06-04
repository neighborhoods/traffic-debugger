FROM golang:1.12-alpine3.9 AS deps
RUN apk add --no-cache git
RUN mkdir /src
WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download

FROM deps AS build-server
RUN mkdir /src/server
COPY ./server /src/server/
RUN go build -o /server ./server/

FROM deps AS build-client
RUN mkdir /src/client
COPY ./client /src/client/
RUN go build -o /client ./client/

FROM alpine:3.9
COPY --from=build-server /server /server
COPY --from=build-client /client /client
RUN chmod +x /server
RUN chmod +x /client

ENTRYPOINT ["/server"]
