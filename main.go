package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/kaspanet/faucet/config"
	"github.com/kaspanet/faucet/database"
	"github.com/kaspanet/faucet/version"
	"github.com/kaspanet/go-secp256k1"
	"github.com/kaspanet/kaspad/dagconfig"
	"github.com/kaspanet/kaspad/txscript"
	"github.com/kaspanet/kaspad/util"
	"github.com/kaspanet/kaspad/util/profiling"
	"github.com/pkg/errors"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kaspanet/kaspad/signal"
	"github.com/kaspanet/kaspad/util/panics"
)

var (
	faucetAddress      util.Address
	faucetPrivateKey   *secp256k1.PrivateKey
	faucetScriptPubKey []byte
)

func main() {
	defer panics.HandlePanic(log, "main", nil)
	interrupt := signal.InterruptListener()

	err := config.Parse()
	if err != nil {
		err := errors.Wrap(err, "Error parsing command-line arguments")
		_, err = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		if err != nil {
			panic(err)
		}
		return
	}

	cfg, err := config.MainConfig()
	if err != nil {
		panic(err)
	}

	// Show version at startup.
	log.Infof("Version %s", version.Version())

	// Enable http profiling server if requested.
	if cfg.Profile != "" {
		profiling.Start(cfg.Profile, log)
	}

	if cfg.Migrate {
		err := database.Migrate(cfg)
		if err != nil {
			panic(errors.Errorf("Error migrating database: %s", err))
		}
		return
	}

	err = database.Connect(cfg)
	if err != nil {
		panic(errors.Errorf("Error connecting to database: %s", err))
	}
	defer func() {
		err := database.Close()
		if err != nil {
			panic(errors.Errorf("Error closing the database: %s", err))
		}
	}()

	privateKeyBytes, err := hex.DecodeString(cfg.PrivateKey)
	if err != nil {
		panic(errors.Wrap(err, "failed to deserialize private key"))
	}

	faucetPrivateKey, _ = secp256k1.DeserializePrivateKeyFromSlice(privateKeyBytes)

	faucetAddress, err = privateKeyToP2PKHAddress(faucetPrivateKey, config.ActiveNetParams())
	if err != nil {
		panic(errors.Errorf("Failed to get P2PKH address from private key: %s", err))
	}

	faucetScriptPubKey, err = txscript.PayToAddrScript(faucetAddress)
	if err != nil {
		panic(errors.Errorf("failed to generate faucetScriptPubKey to address: %s", err))
	}

	shutdownServer := startHTTPServer(cfg.HTTPListen)
	defer shutdownServer()

	<-interrupt
}

// privateKeyToP2PKHAddress generates p2pkh address from private key.
func privateKeyToP2PKHAddress(key *secp256k1.PrivateKey, net *dagconfig.Params) (util.Address, error) {
	publicKey, err := key.SchnorrPublicKey()
	if err != nil {
		return nil, err
	}
	serialized, err := publicKey.SerializeCompressed()
	if err != nil {
		return nil, err
	}
	return util.NewAddressPubKeyHashFromPublicKey(serialized, net.Prefix)
}
