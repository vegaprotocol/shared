package wallet

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	vtypes "code.vegaprotocol.io/vega/core/types"
	"code.vegaprotocol.io/vega/paths"
	commV1 "code.vegaprotocol.io/vega/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/vega/protos/vega/wallet/v1"
)

func TestWalletV2Service_SetupWallet(t *testing.T) {
	t.Skipf("just for manual testing, for now")

	type fields struct {
		networkURL string
		walletName string
		passphrase string
		storePath  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				networkURL: "https://raw.githubusercontent.com/vegaprotocol/networks-internal/main/devnet1/vegawallet-devnet1.toml",
				walletName: vgrand.RandomStr(5),
				passphrase: vgrand.RandomStr(5),
				storePath:  paths.VegaHome,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := NewWalletV2Service(
				tt.fields.walletName,
				tt.fields.passphrase,
				tt.fields.storePath,
				WithNetworkURL(tt.fields.networkURL),
			)
			require.NoError(t, err)

			got, err := w.SetupWallet(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("SetupWallet() got = %v", got)
			}

			tx := &walletpb.SubmitTransactionRequest{
				PubKey:    w.pubKey,
				Propagate: true,
				Command: &walletpb.SubmitTransactionRequest_Transfer{
					Transfer: &commV1.Transfer{
						FromAccountType: vtypes.AccountTypeGeneral,
						To:              "69464e35bcb8e8a2900ca0f87acaf252d50cf2ab2fc73694845a16b7c8a0dc6f",
						ToAccountType:   vtypes.AccountTypeGeneral,
						Asset:           "fc7fd956078fb1fc9db5c19b88f0874c4299b2a7639ad05a47a28c0aef291b55",
						Amount:          "42000000000000",
						Reference:       "Testing the wallet V2 API",
						Kind:            &commV1.Transfer_OneOff{OneOff: &commV1.OneOffTransfer{}},
					},
				},
			}
			signResp, err := w.SignTransaction(context.Background(), tx)
			require.NoError(t, err)

			t.Logf("signResp: %v", signResp)
		})
	}
}
