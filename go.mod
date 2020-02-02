module github.com/kaspanet/faucet

go 1.13

require (
	github.com/golang-migrate/migrate/v4 v4.7.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/jessevdk/go-flags v1.4.0
	github.com/jinzhu/gorm v1.9.11
	github.com/kaspanet/kaspad v0.1.2
	github.com/kaspanet/kasparov v0.1.2
	github.com/pkg/errors v0.9.1
)

replace github.com/kaspanet/kaspad => ../kaspad

replace github.com/kaspanet/kasparov => ../kasparov
