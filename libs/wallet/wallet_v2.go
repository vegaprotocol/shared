package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"go.uber.org/zap"

	"code.vegaprotocol.io/vega/libs/jsonrpc"
	"code.vegaprotocol.io/vega/logging"
	"code.vegaprotocol.io/vega/paths"
	commandspb "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
	"code.vegaprotocol.io/vega/wallet/api"
	"code.vegaprotocol.io/vega/wallet/api/node"
	v1 "code.vegaprotocol.io/vega/wallet/network/store/v1"
	"code.vegaprotocol.io/vega/wallet/wallets"
)

type WalletV2Service struct {
	network        string
	networkFileURL string
	walletName     string
	passphrase     string
	pubKey         string
	nodeAddress    string
	numRetries     uint64

	walletStore  api.WalletStore
	networkStore api.NetworkStore
	log          *logging.Logger

	apiCreateWallet    jsonrpc.Command
	apiDescribeWallet  jsonrpc.Command
	apiListKeys        jsonrpc.Command
	apiGenerateKey     jsonrpc.Command
	apiSendTransaction jsonrpc.Command
}

type WalletV2 interface {
	PublicKey() string
	SendTransaction(ctx context.Context, tx *walletpb.SubmitTransactionRequest) (*commandspb.Transaction, error)
	SendJSONTransaction(ctx context.Context, txJsn string) (*commandspb.Transaction, error)
}

func NewWalletV2Service(log *logging.Logger, config *Config) (*WalletV2Service, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	if config.Name == "" {
		return nil, errors.New("wallet name is required")
	}
	if config.Passphrase == "" {
		return nil, errors.New("passphrase is required")
	}
	if config.StorePath == "" {
		return nil, errors.New("store path is required")
	}

	walletStore, err := wallets.InitialiseStore(config.StorePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	networkStore, err := v1.InitialiseStore(paths.New(config.StorePath))
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise network store: %w", err)
	}

	nodeSelectorBuilder := func(hosts []string, retries uint64) (node.Selector, error) {
		nodes := make([]node.Node, len(hosts))
		for i, host := range hosts {
			nodes[i], err = node.NewRetryingNode(zap.L(), host, retries)
			if err != nil {
				return nil, fmt.Errorf("couldn't initialise retrying node: %w", err)
			}
		}
		return node.NewRoundRobinSelector(zap.L(), nodes...)
	}

	w := &WalletV2Service{
		walletName:     config.Name,
		passphrase:     config.Passphrase,
		walletStore:    walletStore,
		networkStore:   networkStore,
		networkFileURL: config.NetworkFileURL,
		nodeAddress:    config.NodeURL,
		pubKey:         config.PubKey,
		numRetries:     10,
		log:            log.Named("Wallet"),

		apiCreateWallet:    api.NewAdminCreateWallet(walletStore),
		apiDescribeWallet:  api.NewAdminDescribeWallet(walletStore),
		apiListKeys:        api.NewAdminListKeys(walletStore),
		apiGenerateKey:     api.NewAdminGenerateKey(walletStore),
		apiSendTransaction: api.NewAdminSendTransaction(walletStore, networkStore, nodeSelectorBuilder),
	}

	if w.networkFileURL != "" {
		w.network = strings.Split(strings.Split(w.networkFileURL, ".toml")[0], "vegawallet-")[1] // hacky
		networkImportParams := api.AdminImportNetworkParams{
			Name:      w.network,
			URL:       w.networkFileURL,
			Overwrite: true,
		}
		_, errDetails := api.NewAdminImportNetwork(networkStore).Handle(context.Background(), networkImportParams, jsonrpc.RequestMetadata{})
		if errDetails != nil {
			return nil, errors.New(errDetails.Data)
		}
	} else if w.nodeAddress == "" {
		return nil, errors.New("either network file URL or node address must be provided")
	}

	if w.pubKey == "" {
		if err = w.setupWallet(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to login to wallet: %s", err)
		}
	}

	return w, nil
}

func (w *WalletV2Service) PublicKey() string {
	return w.pubKey
}

func (w *WalletV2Service) SendTransaction(ctx context.Context, tx *walletpb.SubmitTransactionRequest) (*commandspb.Transaction, error) {
	if tx.PubKey == "" {
		tx.PubKey = w.pubKey
	}
	tx.Propagate = true
	jsn, err := (&jsonpb.Marshaler{Indent: "  "}).MarshalToString(tx)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal transaction: %w", err)
	}
	return w.SendJSONTransaction(ctx, jsn)
}

func (w *WalletV2Service) SendJSONTransaction(ctx context.Context, payload string) (*commandspb.Transaction, error) {
	txPayload := make(map[string]any)
	if err := json.Unmarshal([]byte(payload), &txPayload); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal transaction payload: %w", err)
	}

	params := api.AdminSendTransactionParams{
		Wallet:      w.walletName,
		Passphrase:  w.passphrase,
		PublicKey:   w.pubKey,
		Retries:     w.numRetries,
		SendingMode: "sync",
		Transaction: txPayload,
	}

	if w.nodeAddress != "" {
		params.NodeAddress = w.nodeAddress
	} else {
		params.Network = w.network
	}

	rawResult, errDetails := w.apiSendTransaction.Handle(ctx, params, jsonrpc.RequestMetadata{})
	if errDetails != nil {
		return nil, errors.New(errDetails.Data)
	}
	return rawResult.(api.AdminSendTransactionResult).Tx, nil
}

func (w *WalletV2Service) setupWallet(ctx context.Context) error {
	_, err := w.describeWallet(ctx)
	if err != nil {
		if err.Error() == api.ErrWalletDoesNotExist.Error() {
			createResp, err := w.createWallet(ctx)
			if err != nil {
				return fmt.Errorf("failed to create wallet: %w", err)
			}
			w.log.Debug("Created wallet",
				logging.String("pubKey", createResp.Key.PublicKey),
				logging.String("wallet", w.walletName),
				logging.String("recoveryPhrase", createResp.Wallet.RecoveryPhrase))
		} else {
			return fmt.Errorf("failed to describe wallet: %w", err)
		}
	} else {
		keysResp, err := w.listKeys(ctx)
		if err != nil {
			return fmt.Errorf("failed to list keys: %w", err)
		}
		if len(keysResp.PublicKeys) == 0 {
			_, err = w.generateKey(ctx)
			if err != nil {
				return fmt.Errorf("failed to generate key: %w", err)
			}
		}
	}

	return nil
}

func (w *WalletV2Service) createWallet(ctx context.Context) (api.AdminCreateWalletResult, error) {
	params := api.AdminCreateWalletParams{
		Wallet:     w.walletName,
		Passphrase: w.passphrase,
	}

	rawResult, errDetails := w.apiCreateWallet.Handle(ctx, params, jsonrpc.RequestMetadata{})
	if errDetails != nil {
		return api.AdminCreateWalletResult{}, errors.New(errDetails.Data)
	}

	result := rawResult.(api.AdminCreateWalletResult)
	w.pubKey = result.Key.PublicKey
	return result, nil
}

func (w *WalletV2Service) describeWallet(ctx context.Context) (api.AdminDescribeWalletResult, error) {
	params := api.AdminDescribeWalletParams{
		Wallet:     w.walletName,
		Passphrase: w.passphrase,
	}

	rawResult, errDetails := w.apiDescribeWallet.Handle(ctx, params, jsonrpc.RequestMetadata{})
	if errDetails != nil {
		return api.AdminDescribeWalletResult{}, errors.New(errDetails.Data)
	}
	return rawResult.(api.AdminDescribeWalletResult), nil
}

func (w *WalletV2Service) listKeys(ctx context.Context) (api.AdminListKeysResult, error) {
	params := api.AdminListKeysParams{
		Wallet:     w.walletName,
		Passphrase: w.passphrase,
	}

	rawResult, errDetails := w.apiListKeys.Handle(ctx, params, jsonrpc.RequestMetadata{})
	if errDetails != nil {
		return api.AdminListKeysResult{}, errors.New(errDetails.Data)
	}

	result := rawResult.(api.AdminListKeysResult)
	if len(result.PublicKeys) == 0 {
		return result, errors.New("no keys found")
	}

	w.pubKey = result.PublicKeys[0].PublicKey
	return result, nil
}

func (w *WalletV2Service) generateKey(ctx context.Context) (api.AdminGenerateKeyResult, error) {
	params := api.AdminGenerateKeyParams{
		Wallet:     w.walletName,
		Passphrase: w.passphrase,
	}

	rawResult, errDetails := w.apiGenerateKey.Handle(ctx, params, jsonrpc.RequestMetadata{})
	if errDetails != nil {
		return api.AdminGenerateKeyResult{}, errors.New(errDetails.Data)
	}

	result := rawResult.(api.AdminGenerateKeyResult)
	w.pubKey = result.PublicKey
	return result, nil
}
