package config

import (
	"time"

	"code.vegaprotocol.io/shared/libs/wallet"
)

type WhaleConfig struct {
	Wallet           *wallet.Config    `yaml:"wallet"`
	OwnerPrivateKeys map[string]string `yaml:"ownerPrivateKeys"`
	TopUpScale       uint64            `yaml:"topUpScale"`
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
