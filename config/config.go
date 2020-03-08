package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/kaspanet/faucet/version"
	"github.com/kaspanet/kaspad/dagconfig"
	"github.com/kaspanet/kaspad/util"
	"github.com/kaspanet/kasparov/logger"
	"github.com/pkg/errors"
)

const (
	defaultLogFilename    = "faucet.log"
	defaultErrLogFilename = "faucet_err.log"
)

var (
	// Default configuration options
	defaultLogDir     = util.AppDataDir("faucet", false)
	defaultDBHost     = "localhost"
	defaultDBPort     = "5432"
	defaultDBSSLMode  = "disable"
	defaultHTTPListen = "0.0.0.0:8081"

	// activeNetParams are the currently active net params
	activeNetParams *dagconfig.Params
)

// Config defines the configuration options for the faucet.
type Config struct {
	ShowVersion bool    `short:"V" long:"version" description:"Display version information and exit"`
	LogDir      string  `long:"logdir" description:"Directory to log output."`
	HTTPListen  string  `long:"listen" description:"HTTP address to listen on (default: 0.0.0.0:8081)"`
	KasparovURL string  `long:"kasparov-url" description:"The kasparov url to connect to"`
	PrivateKey  string  `long:"private-key" description:"Faucet Private key"`
	DBHost      string  `long:"dbhost" description:"Database host"`
	DBPort      string  `long:"dbport" description:"Database port"`
	DBSSLMode   string  `long:"dbsslmode" description:"Database SSL mode {disable, allow, prefer, require, verify-ca, verify-full}"`
	DBUser      string  `long:"dbuser" description:"Database user" required:"true"`
	DBPassword  string  `long:"dbpass" description:"Database password" required:"true"`
	DBName      string  `long:"dbname" description:"Database name" required:"true"`
	Migrate     bool    `long:"migrate" description:"Migrate the database to the latest version. The server will not start when using this flag."`
	FeeRate     float64 `long:"fee-rate" description:"Coins per gram fee rate"`
	TestNet     bool    `long:"testnet" description:"Connect to testnet"`
	SimNet      bool    `long:"simnet" description:"Connect to the simulation test network"`
	DevNet      bool    `long:"devnet" description:"Connect to the development test network"`
}

var cfg *Config

// Parse parses the CLI arguments and returns a config struct.
func Parse() error {
	cfg = &Config{
		LogDir:     defaultLogDir,
		DBHost:     defaultDBHost,
		DBPort:     defaultDBPort,
		DBSSLMode:  defaultDBSSLMode,
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
		if cfg.KasparovURL == "" {
			return errors.New("api-server-url argument is required when --migrate flag is not raised")
		}
		if cfg.PrivateKey == "" {
			return errors.New("private-key argument is required when --migrate flag is not raised")
		}
	}

	err = resolveNetwork(cfg)
	if err != nil {
		return err
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
