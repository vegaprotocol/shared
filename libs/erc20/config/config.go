package config

type TokenConfig struct {
	EthereumAPIAddress   string `yaml:"ethereumAPIAddress"`
	Erc20BridgeAddress   string `yaml:"erc20BridgeAddress"`
	StakingBridgeAddress string `yaml:"stakingBridgeAddress"`
	SyncTimeoutSec       int    `yaml:"syncTimeoutSec"`
}
