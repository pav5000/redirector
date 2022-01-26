FROM golang:alpine AS builder

WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN go build -a -tags netgo -ldflags '-w' -o app github.com/pav5000/redirector/cmd/redirector

FROM scratch

COPY --from=builder /build/app /app
ENTRYPOINT ["/app"]
