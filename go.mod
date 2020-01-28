module github.com/kaspanet/faucet

go 1.13

require (
	github.com/golang-migrate/migrate/v4 v4.7.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/jessevdk/go-flags v1.4.0
	github.com/jinzhu/gorm v1.9.11
	github.com/kaspanet/kaspad v0.1.1-dev
	github.com/kaspanet/kasparov v0.0.0-20200128141254-19ce03a82174
	github.com/pkg/errors v0.8.1
)

replace github.com/kaspanet/kaspad => ../kaspad

replace github.com/kaspanet/kasparov => ../kasparov
