# -- multistage docker build: stage #1: build stage
FROM golang:1.16-alpine AS build

RUN mkdir -p /go/src/github.com/kaspanet/faucet

WORKDIR /go/src/github.com/kaspanet/faucet

RUN apk add --no-cache curl git openssh binutils gcc musl-dev
RUN go get -u golang.org/x/lint/golint

COPY go.mod .
COPY go.sum .

# For development it's useful to have kaspad from filesystem,
# but for deployment we should use the one in github
RUN go mod edit -dropreplace github.com/kaspanet/kaspad 
ARG KASPAD_VERSION
RUN \
    if [ -n "${KASPAD_VERSION}" ] ; then \
        # use replace instead of require - to propagate into nested dependancies
        go mod edit \                       
        -replace github.com/kaspanet/kaspad=github.com/kaspanet/kaspad@$KASPAD_VERSION; \
    fi

RUN go mod download
RUN cp go.mod go.mod.bu
RUN cp go.sum go.sum.bu

COPY . .

# Restore go.mod because `COPY . .` overwrote it
RUN mv go.mod.bu go.mod
RUN mv go.sum.bu go.sum

# Run linters and tests
RUN GOFMT_RESULT=`go fmt ./...`; echo $GOFMT_RESULT; test -z "$GOFMT_RESULT"
RUN go vet ./...
RUN golint -set_exit_status ./...
RUN go test ./...

RUN GOOS=linux go build -a -installsuffix cgo -o faucet .

# --- multistage docker build: stage #2: runtime image
FROM alpine
WORKDIR /app

RUN apk add --no-cache tini

COPY --from=build /go/src/github.com/kaspanet/faucet /app/

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/faucet"]
