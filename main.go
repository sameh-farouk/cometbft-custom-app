package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"

	"test/db"

	cfg "github.com/cometbft/cometbft/config"
	cmtflags "github.com/cometbft/cometbft/libs/cli/flags"
	cmtlog "github.com/cometbft/cometbft/libs/log"
	nm "github.com/cometbft/cometbft/node"
	"github.com/spf13/viper"
)

var (
	homeDir     string
	dbType      string
	dbPath      string
	tbAddresses string
	tbClusterID uint
)

func init() {
	flag.StringVar(&homeDir, "cmt-home", "", "Path to the CometBFT config directory (if empty, uses $HOME/.cometbft)")
	flag.StringVar(&dbType, "db-type", "badger", "Database type: badger, pebble, or tigerbeetle")
	flag.StringVar(&dbPath, "db-path", "", "Path to the database")
	flag.StringVar(&tbAddresses, "tb-addresses", "3000", "TigerBeetle addresses (comma-separated)")
}

func main() {
	flag.Parse()
	if homeDir == "" {
		homeDir = os.ExpandEnv("$HOME/.cometbft")
	}

	// Set defaults from environment if not provided via flags
	if os.Getenv("DB_TYPE") != "" && dbType == "badger" {
		dbType = os.Getenv("DB_TYPE")
	}
	if os.Getenv("DB_PATH") != "" && dbPath == "" {
		dbPath = os.Getenv("DB_PATH")
	}
	if os.Getenv("TB_ADDRESSES") != "" && tbAddresses == "3000" {
		tbAddresses = os.Getenv("TB_ADDRESSES")
	}

	// get ID from environment variable
	// nodeID := os.Getenv("ID")

	config := cfg.DefaultConfig()
	config.SetRoot(homeDir)
	// nodeDir := fmt.Sprintf("node%s", nodeID)
	viper.SetConfigFile(fmt.Sprintf("%s/%s", homeDir, "config/config.toml"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Reading config: %v", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Decoding config: %v", err)
	}
	if err := config.ValidateBasic(); err != nil {
		log.Fatalf("Invalid configuration data: %v", err)
	}

	// Initialize the appropriate database
	var database db.DB
	var err error

	if dbPath == "" {
		dbPath = filepath.Join(homeDir, dbType)
	}

	switch dbType {
	case "badger", "":
		database, err = db.NewBadgerDB(dbPath)
	case "pebble":
		database, err = db.NewPebbleDB(dbPath)
	case "tigerbeetle":
		database, err = db.NewTigerBeetleDBFromMain(tbAddresses)
	default:
		log.Fatalf("Unknown database type: %s", dbType)
	}

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	app := NewKVStoreApplication(database)

	pv := privval.LoadFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
	)

	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		log.Fatalf("failed to load node's key: %v", err)
	}

	logger := cmtlog.NewTMLogger(cmtlog.NewSyncWriter(os.Stdout))
	logger, err = cmtflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel)

	if err != nil {
		log.Fatalf("failed to parse log level: %v", err)
	}

	node, err := nm.NewNode(
		context.Background(),
		config,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		nm.DefaultGenesisDocProviderFunc(config),
		cfg.DefaultDBProvider,
		nm.DefaultMetricsProvider(config.Instrumentation),
		logger,
	)

	if err != nil {
		log.Fatalf("Creating node: %v", err)
	}

	node.Start()
	defer func() {
		node.Stop()
		node.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
