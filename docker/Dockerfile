# -- multistage docker build: stage #1: build stage
FROM golang:1.13-alpine AS build
ARG KASPAD_VERSION=latest
ARG KASPAROV_VERSION=latest

RUN mkdir -p /go/src/github.com/kaspanet/faucet

WORKDIR /go/src/github.com/kaspanet/faucet

RUN apk add --no-cache curl git

COPY go.mod .
COPY go.sum .

# For development it's useful to have kasparov and kaspad from filesystem,
# but for deployment we should use the one in github
RUN go mod edit -dropreplace github.com/kaspanet/kaspad 
RUN go mod edit -dropreplace github.com/kaspanet/kasparov
RUN go mod edit -require github.com/kaspanet/kaspad@$KASPAD_VERSION -require github.com/kaspanet/kasparov@$KASPAROV_VERSION
RUN cp go.mod go.mod.bu

RUN go mod download

COPY . .

# Restore go.mod because `COPY . .` overwrote it
RUN mv go.mod.bu go.mod

# Run linters and tests
RUN GOFMT_RESULT=`go fmt ./...`; echo $GOFMT_RESULT; test -z "$GOFMT_RESULT"
RUN go vet ./...
RUN golint -set_exit_status ./...
RUN go test -gcflags='-l' ./...

RUN  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o faucet .

# --- multistage docker build: stage #2: runtime image
FROM alpine
WORKDIR /app

RUN apk add --no-cache tini

COPY --from=build /go/src/github.com/kaspanet/faucet /app/

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/faucet"]
