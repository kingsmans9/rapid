FROM golang:1.20 as builder

ENV PROJECTPATH=/go/src/github.com/spectrocloud/rapid-agent
ENV LOG_LEVEL DEBUG

WORKDIR $PROJECTPATH

COPY Makefile ./
COPY go.mod ./
COPY go.sum ./
COPY pkg/ ./pkg
COPY cmd/ ./cmd

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

RUN make build
RUN mv ./bin/rapid-agent /rapid-agent

FROM alpine:3.14

COPY --from=builder /rapid-agent /rapid-agent

# non root user cannot read rapid-data directory
# RUN addgroup -S rapid && adduser -S rapid -G rapid
# USER rapid

WORKDIR /

EXPOSE 3000
ENTRYPOINT ["/rapid-agent"]
CMD ["api"]
