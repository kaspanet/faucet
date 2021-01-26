package config

import (
	"fmt"
	"github.com/kaspanet/faucet/logger"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/kaspanet/faucet/version"
	"github.com/kaspanet/kaspad/domain/dagconfig"
	"github.com/kaspanet/kaspad/util"
	"github.com/pkg/errors"
)

const (
	defaultLogFilename    = "faucet.log"
	defaultErrLogFilename = "faucet_err.log"
)

var (
	// Default configuration options
	defaultLogDir     = util.AppDataDir("faucet", false)
	defaultHTTPListen = "0.0.0.0:8081"

	// activeNetParams are the currently active net params
	activeNetParams *dagconfig.Params
)

// Config defines the configuration options for the faucet.
type Config struct {
	ShowVersion bool    `short:"V" long:"version" description:"Display version information and exit"`
	LogDir      string  `long:"logdir" description:"Directory to log output."`
	HTTPListen  string  `long:"listen" description:"HTTP address to listen on default: 0.0.0.0:8081)"`
	RPCServer   string  `long:"rpcserver" short:"s" description:"RPC server to connect to"`
	PrivateKey  string  `long:"private-key" description:"Faucet Private key"`
	DBAddress   string  `long:"dbaddress" description:"Database address" default:"localhost:5432"`
	DBSSLMode   string  `long:"dbsslmode" description:"Database SSL mode" choice:"disable" choice:"allow" choice:"prefer" choice:"require" choice:"verify-ca" choice:"verify-full" default:"disable"`
	DBUser      string  `long:"dbuser" description:"Database user" required:"true"`
	DBPassword  string  `long:"dbpass" description:"Database password" required:"true"`
	DBName      string  `long:"dbname" description:"Database name" required:"true"`
	Migrate     bool    `long:"migrate" description:"Migrate the database to the latest version. The server will not start when using this flag."`
	FeeRate     float64 `long:"fee-rate" description:"Coins per gram fee rate"`
	TestNet     bool    `long:"testnet" description:"Connect to testnet"`
	SimNet      bool    `long:"simnet" description:"Connect to the simulation test network"`
	DevNet      bool    `long:"devnet" description:"Connect to the development test network"`
	Profile     string  `long:"profile" description:"Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536"`
}

var cfg *Config

// Parse parses the CLI arguments and returns a config struct.
func Parse() error {
	cfg = &Config{
		LogDir:     defaultLogDir,
		HTTPListen: defaultHTTPListen,
	}
	parser := flags.NewParser(cfg, flags.HelpFlag)
	_, err := parser.Parse()

	// Show the version and exit if the version flag was specified.
	if cfg.ShowVersion {
		appName := filepath.Base(os.Args[0])
		appName = strings.TrimSuffix(appName, filepath.Ext(appName))
		fmt.Println(appName, "version", version.Version())
		os.Exit(0)
	}

	if err != nil {
		return err
	}

	if !cfg.Migrate {
		if cfg.RPCServer == "" {
			return errors.New("rpcserver argument is required when --migrate flag is not raised")
		}
		if cfg.PrivateKey == "" {
			return errors.New("private-key argument is required when --migrate flag is not raised")
		}
	}

	err = resolveNetwork(cfg)
	if err != nil {
		return err
	}

	if cfg.Profile != "" {
		profilePort, err := strconv.Atoi(cfg.Profile)
		if err != nil || profilePort < 1024 || profilePort > 65535 {
			return errors.New("The profile port must be between 1024 and 65535")
		}
	}

	logFile := filepath.Join(cfg.LogDir, defaultLogFilename)
	errLogFile := filepath.Join(cfg.LogDir, defaultErrLogFilename)
	logger.InitLog(logFile, errLogFile)

	return nil
}

func resolveNetwork(cfg *Config) error {
	// Multiple networks can't be selected simultaneously.
	numNets := 0
	if cfg.TestNet {
		numNets++
	}
	if cfg.SimNet {
		numNets++
	}
	if cfg.DevNet {
		numNets++
	}
	if numNets > 1 {
		return errors.New("multiple net params (testnet, simnet, devnet, etc.) can't be used " +
			"together -- choose one of them")
	}

	activeNetParams = &dagconfig.MainnetParams
	switch {
	case cfg.TestNet:
		activeNetParams = &dagconfig.TestnetParams
	case cfg.SimNet:
		activeNetParams = &dagconfig.SimnetParams
	case cfg.DevNet:
		activeNetParams = &dagconfig.DevnetParams
	}

	return nil
}

// MainConfig is a getter to the main config
func MainConfig() (*Config, error) {
	if cfg == nil {
		return nil, errors.New("No configuration was set for the faucet")
	}
	return cfg, nil
}

// ActiveNetParams returns the currently active net params
func ActiveNetParams() *dagconfig.Params {
	return activeNetParams
}
