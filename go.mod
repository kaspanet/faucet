module github.com/kaspanet/faucet

go 1.13

require (
	github.com/golang-migrate/migrate/v4 v4.7.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/jessevdk/go-flags v1.4.0
	github.com/jinzhu/gorm v1.9.11
	github.com/kaspanet/kaspad v0.0.10-testnet.0.20191217114003-03b7af9a1346
	github.com/kaspanet/kasparov v0.0.0-20191217155301-7d01ee4fab88
	github.com/pkg/errors v0.8.1
)

replace github.com/kaspanet/kaspad => ../kaspad

replace github.com/kaspanet/kasparov => ../kasparov
