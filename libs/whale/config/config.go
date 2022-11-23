package config

import "time"

type WhaleConfig struct {
	WalletPubKey     string            `yaml:"walletPubKey"`
	WalletName       string            `yaml:"walletName"`
	WalletPassphrase string            `yaml:"walletPassphrase"`
	OwnerPrivateKeys map[string]string `yaml:"ownerPrivateKeys"`
	FaucetURL        string            `yaml:"faucetURL"`
	FaucetRateLimit  time.Duration     `yaml:"faucetRateLimit"`
	SyncTimeoutSec   int               `yaml:"syncTimeoutSec"`
	SlackConfig      SlackConfig       `yaml:"slack"`
}

type SlackConfig struct {
	AppToken  string `yaml:"appToken"`
	BotToken  string `yaml:"botToken"`
	ChannelID string `yaml:"channelID"`
	Enabled   bool   `yaml:"enabled"`
}
