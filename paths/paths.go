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

type CachePath string

func (p CachePath) String() string {
	return string(p)
}

var (
	// DataNodeCacheHome is the folder containing the data dedicated to the
	// data-node.
	DataNodeCacheHome = CachePath("data-node")
)

// File structure for configuration
//
// CONFIG_PATH
// 	├── console/
// 	│	├── config.toml
// 	│	└── proxy.toml
// 	├── data-node/
// 	│	└── config.toml
// 	├── faucet/
// 	│	└── config.toml
// 	├── node/
// 	│	├── config.toml
// 	│	└── wallets.toml
// 	├── wallet-cli/
// 	│	└── config.toml
// 	├── wallet-desktop/
// 	│	└── config.toml
// 	└── wallet-service/
// 		└── networks/

type ConfigPath string

func (p ConfigPath) String() string {
	return string(p)
}

var (
	// ConsoleConfigHome is the folder containing the configuration files
	// dedicated to the console.
	ConsoleConfigHome = ConfigPath("console")

	// ConsoleDefaultConfigFile is the default configuration file for the
	// console.
	ConsoleDefaultConfigFile = ConfigPath(filepath.Join(ConsoleConfigHome.String(), "config.toml"))

	// ConsoleProxyConfigFile is the configuration file for the
	// console proxy.
	ConsoleProxyConfigFile = ConfigPath(filepath.Join(ConsoleConfigHome.String(), "proxy.toml"))

	// DataNodeConfigHome is the folder containing the configuration files
	// dedicated to the node.
	DataNodeConfigHome = ConfigPath("data-node")

	// DataNodeDefaultConfigFile is the default configuration file for the
	// data-node.
	DataNodeDefaultConfigFile = ConfigPath(filepath.Join(DataNodeConfigHome.String(), "config.toml"))

	// FaucetConfigHome is the folder containing the configuration files
	// dedicated to the node.
	FaucetConfigHome = ConfigPath("faucet")

	// FaucetDefaultConfigFile is the default configuration file for the
	// data-node.
	FaucetDefaultConfigFile = ConfigPath(filepath.Join(FaucetConfigHome.String(), "config.toml"))

	// NodeConfigHome is the folder containing the configuration files dedicated
	// to the node.
	NodeConfigHome = ConfigPath("node")

	// NodeDefaultConfigFile is the default configuration file for the node.
	NodeDefaultConfigFile = ConfigPath(filepath.Join(NodeConfigHome.String(), "config.toml"))

	// NodeWalletsConfigFile is the configuration file for the node wallets.
	NodeWalletsConfigFile = ConfigPath(filepath.Join(NodeConfigHome.String(), "wallets.encrypted"))

	// WalletCLIConfigHome is the folder containing the configuration files
	// dedicated to the wallet CLI.
	WalletCLIConfigHome = ConfigPath("wallet-cli")

	// WalletCLIDefaultConfigFile is the default configuration file for the
	// wallet CLI.
	WalletCLIDefaultConfigFile = ConfigPath(filepath.Join(WalletCLIConfigHome.String(), "config.toml"))

	// WalletDesktopConfigHome is the folder containing the configuration files
	// dedicated to the wallet desktop application.
	WalletDesktopConfigHome = ConfigPath("wallet-desktop")

	// WalletDesktopDefaultConfigFile is the default configuration file for the
	// wallet desktop application.
	WalletDesktopDefaultConfigFile = ConfigPath(filepath.Join(WalletDesktopConfigHome.String(), "config.toml"))

	// WalletServiceConfigHome is the folder containing the configuration files
	// dedicated to the wallet service application.
	WalletServiceConfigHome = ConfigPath("wallet-service")

	// WalletServiceNetworksConfigHome is the folder containing the
	// configuration files dedicated to the networks.
	WalletServiceNetworksConfigHome = ConfigPath(filepath.Join(WalletServiceConfigHome.String(), "networks"))
)

// File structure for data
//
// DATA_PATH
// 	├── node/
// 	│	└── wallets/
// 	│		├── vega/
// 	│		│	└── vega.timestamp
// 	│		└── ethereum/
// 	│			└── eth-node-wallet
// 	├── faucet/
// 	│	└── wallets/
// 	│		└── vega.timestamp
// 	├── wallets/
// 	│	├── vega-wallet-1
// 	│	└── vega-wallet-2
// 	└── wallet-service/
// 		└── rsa-keys/
// 			├── private.pem
// 			└── public.pem

type DataPath string

func (p DataPath) String() string {
	return string(p)
}

var (
	// NodeDataHome is the folder containing the data dedicated to the node.
	NodeDataHome = DataPath("node")

	// NodeWalletsDataHome is the folder containing the data dedicated to the
	// node wallets.
	NodeWalletsDataHome = DataPath(filepath.Join(NodeDataHome.String(), "wallets"))

	// VegaNodeWalletsDataHome is the folder containing the vega wallet
	// dedicated to the node.
	VegaNodeWalletsDataHome = DataPath(filepath.Join(NodeWalletsDataHome.String(), "vega"))

	// EthereumNodeWalletsDataHome is the folder containing the ethereum wallet
	// dedicated to the node.
	EthereumNodeWalletsDataHome = DataPath(filepath.Join(NodeWalletsDataHome.String(), "ethereum"))

	// FaucetDataHome is the folder containing the data dedicated to the faucet.
	FaucetDataHome = DataPath("faucet")

	// FaucetWalletsDataHome is the folder containing the data dedicated to the
	// faucet wallets.
	FaucetWalletsDataHome = DataPath(filepath.Join(FaucetDataHome.String(), "wallets"))

	// WalletsDataHome is the folder containing the user wallets.
	WalletsDataHome = DataPath("wallets")

	// WalletServiceDataHome is the folder containing the data dedicated to the
	// wallet service.
	WalletServiceDataHome = DataPath("wallet-service")

	// WalletServiceRSAKeysDataHome is the folder containing the RSA keys used by
	// the wallet service.
	WalletServiceRSAKeysDataHome = DataPath(filepath.Join(WalletServiceDataHome.String(), "rsa-keys"))

	// WalletServicePublicRSAKeyDataFile is the file containing the public RSA key
	// used by the wallet service.
	WalletServicePublicRSAKeyDataFile = DataPath(filepath.Join(WalletServiceRSAKeysDataHome.String(), "public.pem"))

	// WalletServicePrivateRSAKeyDataFile is the file containing the private RSA key
	// used by the wallet service.
	WalletServicePrivateRSAKeyDataFile = DataPath(filepath.Join(WalletServiceRSAKeysDataHome.String(), "private.pem"))
)

// File structure for state
//
// STATE_HOME
// 	├── data-node/
// 	│	├── logs/
// 	│	└── storage/
// 	├── node/
// 	│	├── logs/
// 	│	├── checkpoints/
// 	│	└── snapshots/
// 			└── ldb
// 	├── wallet-cli/
// 	│	└── logs/
// 	├── wallet-desktop/
// 	│	└── logs/
// 	└── wallet-service/
// 		└── logs/

type StatePath string

func (p StatePath) String() string {
	return string(p)
}

var (
	// DataNodeStateHome is the folder containing the state dedicated to the
	// data-node.
	DataNodeStateHome = StatePath("data-node")

	// DataNodeLogsHome is the folder containing the logs of the data-node.
	DataNodeLogsHome = StatePath(filepath.Join(DataNodeStateHome.String(), "logs"))

	// DataNodeStorageHome is the folder containing the data storage of the
	// data-node.
	DataNodeStorageHome = StatePath(filepath.Join(DataNodeStateHome.String(), "storage"))

	// NodeStateHome is the folder containing the state of the node.
	NodeStateHome = StatePath("node")

	// NodeLogsHome is the folder containing the logs of the node.
	NodeLogsHome = StatePath(filepath.Join(NodeStateHome.String(), "logs"))

	// CheckpointStateHome is the folder containing the checkpoint files
	// of to the node.
	CheckpointStateHome = StatePath(filepath.Join(NodeStateHome.String(), "checkpoints"))

	// SnapshotStateHome is the folder containing the snapshot files
	// of to the node.
	SnapshotStateHome = StatePath(filepath.Join(NodeStateHome.String(), "snapshots"))

	// SnapshotDBStateFile is the DB file for GoLevelDB used in snapshots
	SnapshotDBStateFile = StatePath(filepath.Join(SnapshotStateHome.String(), "ldb"))

	// WalletCLIStateHome is the folder containing the state of the wallet CLI.
	WalletCLIStateHome = StatePath("wallet-cli")

	// WalletCLILogsHome is the folder containing the logs of the wallet CLI.
	WalletCLILogsHome = StatePath(filepath.Join(WalletCLIStateHome.String(), "logs"))

	// WalletDesktopStateHome is the folder containing the state of the wallet
	// desktop.
	WalletDesktopStateHome = StatePath("wallet-desktop")

	// WalletDesktopLogsHome is the folder containing the logs of the wallet
	// desktop.
	WalletDesktopLogsHome = StatePath(filepath.Join(WalletDesktopStateHome.String(), "logs"))

	// WalletServiceStateHome is the folder containing the state of the node.
	WalletServiceStateHome = StatePath("wallet-service")

	// WalletServiceLogsHome is the folder containing the logs of the node.
	WalletServiceLogsHome = StatePath(filepath.Join(WalletServiceStateHome.String(), "logs"))
)
