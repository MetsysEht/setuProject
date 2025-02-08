# Stage 1 - Generate proto
###########################
FROM golang:alpine AS proto-builder

WORKDIR /src

ENV GO111MODULE=auto
RUN go env -w GO111MODULE=auto

ADD Makefile /src
ADD proto/ /src/proto/
ADD buf.yaml buf.gen.yaml /src/

RUN apk add --no-cache make
RUN make proto-deps proto-refresh

# Stage 2 - Compilation build stage
######################################
FROM golang:alpine AS builder

ENV CGO_ENABLED 0
RUN mkdir /app
WORKDIR /app
RUN apk add --no-cache make

COPY --from=proto-builder /src/rpc /app/rpc

ADD . /app/

ARG TARGETARCH
ARG TARGETOS

RUN make go-build-api GOARCH=$TARGETARCH GOOS=$TARGETOS

# Stage 3 - Binary build stage
######################################
FROM golang:alpine

COPY --from=builder /app/bin/api /app/
COPY --from=builder /app/config/ /app/config/
COPY entrypoint.sh /app/

ENV WORKDIR=/app
ENV DUMB_INIT_SETSID=0
WORKDIR /app

RUN apk add --update --no-cache dumb-init su-exec curl tzdata

RUN chmod +x entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh", "api"]
