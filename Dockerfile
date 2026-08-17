FROM golang:1.23.12 AS builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test ./... \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go vet ./... \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o /out/transitmanifest ./cmd/transitmanifest

FROM scratch
COPY --from=builder /out/transitmanifest /transitmanifest
ENTRYPOINT ["/transitmanifest"]

