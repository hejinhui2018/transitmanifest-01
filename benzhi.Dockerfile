ARG GO_IMAGE=golang:1.23.12
FROM ${GO_IMAGE}

ENV GOTOOLCHAIN=local \
    CGO_ENABLED=0

WORKDIR /workspace
COPY go.mod ./
RUN GOPROXY=https://goproxy.cn,direct go mod download
COPY . .
RUN go test ./... -count=1
RUN go vet ./...
RUN go build ./...

CMD ["go", "test", "./...", "-count=1"]
