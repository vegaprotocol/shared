package paths

const (
	// LongestPathNameLen is the length of the longest path name. It is used
	// for text formatting.
	LongestPathNameLen = 35
)

type ListPathsResponse struct {
	CachePaths  map[string]string `json:"cachePaths"`
	ConfigPaths map[string]string `json:"configPaths"`
	DataPaths   map[string]string `json:"dataPaths"`
	StatePaths  map[string]string `json:"statePaths"`
}

func List(vegaPaths Paths) *ListPathsResponse {
	return &ListPathsResponse{
		CachePaths: map[string]string{
			"DataNodeCacheHome": vegaPaths.CachePathFor(DataNodeCacheHome),
		},
		ConfigPaths: map[string]string{
			"ConsoleConfigHome":               vegaPaths.ConfigPathFor(ConsoleConfigHome),
			"ConsoleDefaultConfigFile":        vegaPaths.ConfigPathFor(ConsoleDefaultConfigFile),
			"ConsoleProxyConfigFile":          vegaPaths.ConfigPathFor(ConsoleProxyConfigFile),
			"DataNodeConfigHome":              vegaPaths.ConfigPathFor(DataNodeConfigHome),
			"DataNodeDefaultConfigFile":       vegaPaths.ConfigPathFor(DataNodeDefaultConfigFile),
			"FaucetConfigHome":                vegaPaths.ConfigPathFor(FaucetConfigHome),
			"FaucetDefaultConfigFile":         vegaPaths.ConfigPathFor(FaucetDefaultConfigFile),
			"NodeConfigHome":                  vegaPaths.ConfigPathFor(NodeConfigHome),
			"NodeDefaultConfigFile":           vegaPaths.ConfigPathFor(NodeDefaultConfigFile),
			"NodeWalletsConfigFile":           vegaPaths.ConfigPathFor(NodeWalletsConfigFile),
			"WalletCLIConfigHome":             vegaPaths.ConfigPathFor(WalletCLIConfigHome),
			"WalletCLIDefaultConfigFile":      vegaPaths.ConfigPathFor(WalletCLIDefaultConfigFile),
			"WalletDesktopConfigHome":         vegaPaths.ConfigPathFor(WalletDesktopConfigHome),
			"WalletDesktopDefaultConfigFile":  vegaPaths.ConfigPathFor(WalletDesktopDefaultConfigFile),
			"WalletServiceConfigHome":         vegaPaths.ConfigPathFor(WalletServiceConfigHome),
			"WalletServiceNetworksConfigHome": vegaPaths.ConfigPathFor(WalletServiceNetworksConfigHome),
		},
		DataPaths: map[string]string{
			"NodeDataHome":                       vegaPaths.DataPathFor(NodeDataHome),
			"NodeWalletsDataHome":                vegaPaths.DataPathFor(NodeWalletsDataHome),
			"VegaNodeWalletsDataHome":            vegaPaths.DataPathFor(VegaNodeWalletsDataHome),
			"EthereumNodeWalletsDataHome":        vegaPaths.DataPathFor(EthereumNodeWalletsDataHome),
			"FaucetDataHome":                     vegaPaths.DataPathFor(FaucetDataHome),
			"FaucetWalletsDataHome":              vegaPaths.DataPathFor(FaucetWalletsDataHome),
			"WalletsDataHome":                    vegaPaths.DataPathFor(WalletsDataHome),
			"WalletServiceDataHome":              vegaPaths.DataPathFor(WalletServiceDataHome),
			"WalletServiceRSAKeysDataHome":       vegaPaths.DataPathFor(WalletServiceRSAKeysDataHome),
			"WalletServicePublicRSAKeyDataFile":  vegaPaths.DataPathFor(WalletServicePublicRSAKeyDataFile),
			"WalletServicePrivateRSAKeyDataFile": vegaPaths.DataPathFor(WalletServicePrivateRSAKeyDataFile),
		},
		StatePaths: map[string]string{
			"DataNodeStateHome":      vegaPaths.StatePathFor(DataNodeStateHome),
			"DataNodeLogsHome":       vegaPaths.StatePathFor(DataNodeLogsHome),
			"DataNodeStorageHome":    vegaPaths.StatePathFor(DataNodeStorageHome),
			"NodeStateHome":          vegaPaths.StatePathFor(NodeStateHome),
			"NodeLogsHome":           vegaPaths.StatePathFor(NodeLogsHome),
			"CheckpointStateHome":    vegaPaths.StatePathFor(CheckpointStateHome),
			"SnapshotStateHome":      vegaPaths.StatePathFor(SnapshotStateHome),
			"SnapshotDBStateFile":    vegaPaths.StatePathFor(SnapshotDBStateFile),
			"WalletCLIStateHome":     vegaPaths.StatePathFor(WalletCLIStateHome),
			"WalletCLILogsHome":      vegaPaths.StatePathFor(WalletCLILogsHome),
			"WalletDesktopStateHome": vegaPaths.StatePathFor(WalletDesktopStateHome),
			"WalletDesktopLogsHome":  vegaPaths.StatePathFor(WalletDesktopLogsHome),
			"WalletServiceStateHome": vegaPaths.StatePathFor(WalletServiceStateHome),
			"WalletServiceLogsHome":  vegaPaths.StatePathFor(WalletServiceLogsHome),
		},
	}
}
