FROM golang:alpine AS builder

WORKDIR /src/

ENV GO111MODULE=auto
RUN go env -w GO111MODULE=auto

COPY go.mod .
COPY go.sum .

# Install reflex
RUN go install github.com/cespare/reflex@latest

# Make sure $GOPATH/bin is in your PATH
ENV PATH=$PATH:/go/bin

RUN go mod download

COPY . .

RUN go build -o ./run ./cmd/api

FROM alpine:latest AS app
RUN apk --no-cache add ca-certificates
WORKDIR /root/

#Copy executable from builder
COPY --from=builder /src/run .

EXPOSE 8080
CMD ["./run"]