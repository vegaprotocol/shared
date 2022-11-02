package httpcore

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"code.vegaprotocol.io/shared/libs/errors"
)

// Header represents one key-value pair.
type Header struct {
	Key, Value string
}

// DoHTTP does an HTTP call, checks for standard failure cases, and
// returns the body contents.
func DoHTTP(
	ctx context.Context,
	cli *http.Client, url *url.URL, method string, body io.Reader,
	headers []Header,
) ([]byte, error) {
	req, _ := http.NewRequestWithContext(ctx, method, url.String(), body)
	for _, hdr := range headers {
		req.Header.Add(hdr.Key, hdr.Value)
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.ErrServerResponseNone
	}
	if resp.Body == nil {
		return nil, errors.ErrServerResponseEmpty
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.ErrServerResponseReadFail
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad HTTP status code: %d (%s)", resp.StatusCode, content)
	}
	return content, err
}

// NewHTTPClient creates an http.Client with default parameters.
func NewHTTPClient() *http.Client {
	dialContext := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	transport := http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialContext.DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{Transport: &transport, Timeout: time.Second * 30}
}
