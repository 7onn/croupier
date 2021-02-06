FROM golang:1.15.7 as builder
ARG TAG
ENV GOPATH=/go
ENV GOCACHE=/go/src/github.com/devbytom/croupier/src/build-cache

WORKDIR /go/src/github.com/devbytom/croupier/src/
COPY . .

RUN go mod download && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --ldflags '-X main.VERSION=$(TAG) -w -extldflags "-static"' -tags netgo -o server ./

FROM centurylink/ca-certs
COPY --from=builder /go/src/github.com/devbytom/croupier/src/server /bin/server
CMD ["/bin/server"]
