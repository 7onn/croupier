FROM golang:1.17.7 as builder
ARG TAG=latest
ENV GOPATH=/go
ENV GOCACHE=/go/src/github.com/7onn/croupier/app/croupier-api/build-cache

WORKDIR /go/src/github.com/7onn/croupier
COPY . .

RUN go mod download && \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  go build \
    -a \
    --ldflags '-X main.VERSION=$(TAG) -w -extldflags "-static"' \
    -tags netgo \
    -o server ./app/croupier-api

FROM centurylink/ca-certs
COPY --from=builder /go/src/github.com/7onn/croupier/server /bin/server
ENTRYPOINT ["/bin/server"]
