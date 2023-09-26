FROM golang:1.21.1-alpine3.18 AS builder

WORKDIR /build

RUN echo "nobody:x:1337:1337:::" > etcpasswd
RUN echo "nobody::1337:nobody" > etcgroup

COPY go.mod go.sum .
RUN go mod download

COPY *.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tun .

FROM scratch

COPY --from=builder /build/etcpasswd /etc/passwd
COPY --from=builder /build/etcgroup /etc/group
COPY --from=builder --chown=nobody:nobody /build/tun /tun

WORKDIR /
USER nobody

ENTRYPOINT [ "/tun" ]