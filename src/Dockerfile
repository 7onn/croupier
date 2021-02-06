FROM golang:1.15.7 as builder
ARG TAG=latest
ENV GOPATH=/go
ENV GOPATH=/go/src/github.com/devbytom/croupier/server/build-go-path
ENV GOCACHE=/go/src/github.com/devbytom/croupier/server/build-cache
WORKDIR /go/src/github.com/devbytom/croupier/server/
COPY . .
RUN go mod download && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --ldflags '-X main.VERSION=$(TAG) -w -extldflags "-static"' -tags netgo -o server ./

FROM centurylink/ca-certs
COPY --from=builder /go/src/github.com/devbytom/croupier/server /server
ENTRYPOINT ["/server"]
