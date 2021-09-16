package paths

import "path/filepath"

var (
	// VegaHome is the name of the Vega folder for every type of file structure.
	VegaHome = "vega"
)

// File structure for cache
//
// CACHE_PATH
// 	└── data-node/

var (
	// DataNodeCacheHome is the folder containing the data dedicated to the
	// data-node.
	DataNodeCacheHome = "data-node"
)

// File structure for configuration
//
// CONFIG_PATH
// 	├── console/
// 	│	├── config.toml
// 	│	└── proxy.toml
// 	├── data-node/
// 	│	└── config.toml
// 	├── node/
// 	│	├── config.toml
// 	│	└── wallets.toml
// 	├── wallet-cli/
// 	│	└── config.toml
// 	├── wallet-desktop/
// 	│	└── config.toml
// 	└── wallet-service/
// 		└── config.toml

var (
	// ConsoleConfigHome is the folder containing the configuration files
	// dedicated to the console.
	ConsoleConfigHome = "console"

	// ConsoleDefaultConfigFile is the default configuration file for the
	// console.
	ConsoleDefaultConfigFile = filepath.Join(ConsoleConfigHome, "config.toml")

	// ConsoleProxyConfigFile is the configuration file for the
	// console proxy.
	ConsoleProxyConfigFile = filepath.Join(ConsoleConfigHome, "proxy.toml")

	// DataNodeConfigHome is the folder containing the configuration files
	// dedicated to the node.
	DataNodeConfigHome = "data-node"

	// DataNodeDefaultConfigFile is the default configuration file for the
	// data-node.
	DataNodeDefaultConfigFile = filepath.Join(DataNodeConfigHome, "config.toml")

	// NodeConfigHome is the folder containing the configuration files dedicated
	// to the node.
	NodeConfigHome = "node"

	// NodeDefaultConfigFile is the default configuration file for the node.
	NodeDefaultConfigFile = filepath.Join(NodeConfigHome, "config.toml")

	// NodeWalletsConfigFile is the configuration file for the node wallets.
	NodeWalletsConfigFile = filepath.Join(NodeConfigHome, "wallets.encrypted")

	// WalletCLIConfigHome is the folder containing the configuration files
	// dedicated to the wallet CLI.
	WalletCLIConfigHome = "wallet-cli"

	// WalletCLIDefaultConfigFile is the default configuration file for the
	// wallet CLI.
	WalletCLIDefaultConfigFile = filepath.Join(WalletCLIConfigHome, "config.toml")

	// WalletDesktopConfigHome is the folder containing the configuration files
	// dedicated to the wallet desktop application.
	WalletDesktopConfigHome = "wallet-desktop"

	// WalletDesktopDefaultConfigFile is the default configuration file for the
	// wallet desktop application.
	WalletDesktopDefaultConfigFile = filepath.Join(WalletDesktopConfigHome, "config.toml")

	// WalletServiceConfigHome is the folder containing the configuration files
	// dedicated to the wallet desktop application.
	WalletServiceConfigHome = "wallet-service"

	// WalletServiceDefaultConfigFile is the default configuration file for the
	// wallet desktop application.
	WalletServiceDefaultConfigFile = filepath.Join(WalletServiceConfigHome, "config.toml")
)

// File structure for data
//
// DATA_PATH
// 	├── node/
// 	│	└── wallets/
// 	│		├── vega/
// 	│		│	└── vega-node-wallet
// 	│		└── ethereum/
// 	│			└── eth-node-wallet
// 	├── wallets/
// 	│	├── vega-wallet-1
// 	│	└── vega-wallet-2
// 	└── wallet-service/
// 		└── rsa-keys/
// 			├── private.pem
// 			└── public.pem

var (
	// NodeDataHome is the folder containing the data dedicated to the node.
	NodeDataHome = "node"

	// NodeWalletsDataHome is the folder containing the data dedicated to the
	// node wallets.
	NodeWalletsDataHome = filepath.Join(NodeDataHome, "wallets")

	// VegaNodeWalletsDataHome is the folder containing the vega wallet
	// dedicated to the node.
	VegaNodeWalletsDataHome = filepath.Join(NodeWalletsDataHome, "vega")

	// EthereumNodeWalletsDataHome is the folder containing the ethereum wallet
	// dedicated to the node.
	EthereumNodeWalletsDataHome = filepath.Join(NodeWalletsDataHome, "ethereum")

	// WalletsDataHome is the folder containing the user wallets.
	WalletsDataHome = "wallets"

	// WalletServiceDataHome is the folder containing the data dedicated to the
	// wallet service.
	WalletServiceDataHome = "wallet-service"

	// WalletServiceRSAKeysDataHome is the folder containing the RSA keys used by
	// the wallet service.
	WalletServiceRSAKeysDataHome = filepath.Join(WalletServiceDataHome, "rsa-keys")

	// WalletServicePublicRSAKeyDataFile is the file containing the public RSA key
	// used by the wallet service.
	WalletServicePublicRSAKeyDataFile = filepath.Join(WalletServiceRSAKeysDataHome, "public.pem")

	// WalletServicePrivateRSAKeyDataFile is the file containing the private RSA key
	// used by the wallet service.
	WalletServicePrivateRSAKeyDataFile = filepath.Join(WalletServiceRSAKeysDataHome, "private.pem")
)

// File structure for state
//
// STATE_HOME
// 	├── node/
// 	│	├── logs/
// 	│	├── checkpoints/
// 	│	└── snapshots/
// 	├── wallet-cli/
// 	│	└── logs/
// 	├── wallet-desktop/
// 	│	└── logs/
// 	└── wallet-service/
// 		└── logs/

var (
	// DataNodeStateHome is the folder containing the state dedicated to the
	// data-node.
	DataNodeStateHome = "data-node"

	// DataNodeLogsHome is the folder containing the logs of the data-node.
	DataNodeLogsHome = filepath.Join(DataNodeStateHome, "logs")

	// NodeStateHome is the folder containing the state of the node.
	NodeStateHome = "node"

	// NodeLogsHome is the folder containing the logs of the node.
	NodeLogsHome = filepath.Join(NodeStateHome, "logs")

	// CheckpointStateHome is the folder containing the checkpoint files
	// of to the node.
	CheckpointStateHome = filepath.Join(NodeStateHome, "checkpoints")

	// SnapshotStateHome is the folder containing the snapshot files
	// of to the node.
	SnapshotStateHome = filepath.Join(NodeStateHome, "snapshots")

	// WalletCLIStateHome is the folder containing the state of the wallet CLI.
	WalletCLIStateHome = "wallet-cli"

	// WalletCLILogsHome is the folder containing the logs of the wallet CLI.
	WalletCLILogsHome = filepath.Join(WalletCLIStateHome, "logs")

	// WalletDesktopStateHome is the folder containing the state of the wallet
	// desktop.
	WalletDesktopStateHome = "wallet-desktop"

	// WalletDesktopLogsHome is the folder containing the logs of the wallet
	// desktop.
	WalletDesktopLogsHome = filepath.Join(WalletDesktopStateHome, "logs")

	// WalletServiceStateHome is the folder containing the state of the node.
	WalletServiceStateHome = "wallet-service"

	// WalletServiceLogsHome is the folder containing the logs of the node.
	WalletServiceLogsHome = filepath.Join(WalletServiceStateHome, "logs")
)
