FROM golang:1.20 as builder

ENV PROJECTPATH=/go/src/github.com/pavansokkenagaraj/rapid-agent
ENV LOG_LEVEL DEBUG

WORKDIR $PROJECTPATH

COPY Makefile ./
COPY go.mod ./
COPY go.sum ./
COPY pkg/ ./pkg
COPY cmd/ ./cmd

RUN make build
RUN mv ./bin/ra /rapid-agent

#  TODO: PA1 - use multi-stage build to reduce image size, change owner of binary to non-root user 
#        and run as non-root user in production environment to improve security

EXPOSE 3000
ENTRYPOINT ["/rapid-agent"]
CMD ["api"]
