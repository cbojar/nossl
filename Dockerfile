FROM golang:1 AS build

ENV CGO_ENABLED=0

RUN update-ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build ./... && go build -o nossl


FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

COPY --from=build --chown=1001 /build/nossl ./nossl

USER 1001

ENTRYPOINT ["/app/nossl"]
