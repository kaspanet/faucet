module github.com/kaspanet/faucet

go 1.16

require (
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-pg/pg/v9 v9.2.1
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/kaspanet/go-secp256k1 v0.0.7
	github.com/kaspanet/kaspad v0.12.0
	github.com/lib/pq v1.10.6 // indirect
	github.com/pkg/errors v0.9.1
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/encoding v0.3.5 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/crypto v0.0.0-20220517005047-85d78b3ac167 // indirect
	golang.org/x/net v0.0.0-20220516155154-20f960328961 // indirect
	golang.org/x/sys v0.0.0-20220513210249-45d2b4557a2a // indirect
	google.golang.org/genproto v0.0.0-20220505152158-f39f71e6c8f3 // indirect
	google.golang.org/grpc v1.46.2 // indirect
)

replace github.com/kaspanet/kaspad => ../kaspad
