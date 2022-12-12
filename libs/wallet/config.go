package wallet

// Config describes the settings for running an internal wallet server.
type Config struct {
	Name           string `yaml:"name"`
	Passphrase     string `yaml:"passphrase"`
	PubKey         string `yaml:"pubKey"`
	StorePath      string `yaml:"storePath"`
	NetworkFileURL string `yaml:"networkFileURL"`
	NodeURL        string `yaml:"nodeURL"`
}
