FROM golang:1.21 AS builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o almanac .

FROM scratch
WORKDIR /config
COPY --from=builder /build/almanac /almanac
ENTRYPOINT ["/almanac", "serve", "--content-dir", "/content"]
