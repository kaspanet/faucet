# -- multistage docker build: stage #1: build stage
FROM golang:1.14-alpine AS build

RUN mkdir -p /go/src/github.com/kaspanet/faucet

WORKDIR /go/src/github.com/kaspanet/faucet

RUN apk add --no-cache curl git openssh binutils gcc musl-dev
RUN go get -u golang.org/x/lint/golint

COPY go.mod .
COPY go.sum .

# For development it's useful to have kasparov and kaspad from filesystem,
# but for deployment we should use the one in github
RUN go mod edit -dropreplace github.com/kaspanet/kaspad 
RUN go mod edit -dropreplace github.com/kaspanet/kasparov
ARG KASPAD_VERSION
ARG KASPAROV_VERSION
RUN \
    if [ -n "${KASPAD_VERSION}" ] || [ -n "${KASPAROV_VERSION}" ] ; then \
        if [ -z "${KASPAD_VERSION}" ] || [ -z "${KASPAROV_VERSION}" ] ; then \
            >&2 echo "KASPAD_VERSION and KASPAROV_VERSION must be both set or both unset"; \
            exit 1; \
        fi; \
        # use replace instead of require - to propagate into nested dependancies
        go mod edit \                       
        -replace github.com/kaspanet/kaspad=github.com/kaspanet/kaspad@$KASPAD_VERSION \
        -replace github.com/kaspanet/kasparov=github.com/kaspanet/kasparov@$KASPAROV_VERSION; \
    fi

RUN cp go.mod go.mod.bu
RUN go mod download

COPY . .

# Restore go.mod because `COPY . .` overwrote it
RUN mv go.mod.bu go.mod

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
