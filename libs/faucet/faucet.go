package faucet

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"code.vegaprotocol.io/shared/libs/errors"
	"code.vegaprotocol.io/shared/libs/httpcore"
)

const (
	urlMint = "/api/v1/mint"
)

// FaucetClient stores state for a faucet client which communicates with REST.
type client struct {
	Address url.URL

	httpClient *http.Client
}

type mintRequest struct {
	Party  string `json:"party"`
	Amount string `json:"amount"` // uint256
	Asset  string `json:"asset"`
}

type mintResponse struct {
	Success bool `json:"success"`
}

// New returns a new faucet client.
func New(addr url.URL) *client {
	cli := client{
		Address: addr,

		httpClient: httpcore.NewHTTPClient(),
	}

	return &cli
}

// NilFaucetClient returns nil. Used in testing only.
func NilFaucetClient() (*client, error) {
	var cli *client
	return cli, nil
}

// GetAddress gets the address of the node.
func (c *client) GetAddress() (url.URL, error) {
	if c == nil {
		return url.URL{}, errors.ErrNil
	}
	return c.Address, nil
}

func (c *client) Mint(ctx context.Context, amount string, asset, party string) (bool, error) {
	if c == nil {
		return false, errors.ErrNil
	}

	req := mintRequest{
		Amount: amount,
		Asset:  asset,
		Party:  party,
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return false, err
	}

	content, err := httpcore.DoHTTP(ctx, c.httpClient, c.url(urlMint), http.MethodPost, bytes.NewBuffer(payload), nil)
	if err != nil {
		return false, err
	}

	resp := mintResponse{}
	if err := json.Unmarshal(content, &resp); err != nil {
		return false, err
	}

	return resp.Success, nil
}

func (c *client) url(path string) *url.URL {
	return c.Address.ResolveReference(&url.URL{Path: path})
}
