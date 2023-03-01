package account

import (
	"context"
	"errors"
	"testing"

	"code.vegaprotocol.io/shared/libs/num"
	"code.vegaprotocol.io/shared/libs/types"
	"code.vegaprotocol.io/vega/protos/vega"
)

func TestService_topUp(t *testing.T) {
	type fields struct {
		name         string
		pubKey       string
		coinProvider CoinProvider
	}
	type args struct {
		ctx       context.Context
		asset     *vega.Asset
		askAmount *num.Uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				name:         "testbot",
				pubKey:       "something",
				coinProvider: mockCoinProvider{err: nil},
			},
			args: args{
				ctx:       context.Background(),
				asset:     &vega.Asset{Id: "fBTC"},
				askAmount: num.NewUint(100),
			},
			wantErr: false,
		}, {
			name: "error",
			fields: fields{
				name:         "testbot",
				pubKey:       "something",
				coinProvider: mockCoinProvider{err: errors.New("failed to top up")},
			},
			args: args{
				ctx:       context.Background(),
				asset:     &vega.Asset{Id: "fBTC"},
				askAmount: num.NewUint(100),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Service{
				name:         tt.fields.name,
				pubKey:       tt.fields.pubKey,
				coinProvider: tt.fields.coinProvider,
			}
			if err := a.topUp(tt.args.ctx, tt.args.asset, tt.args.askAmount); (err != nil) != tt.wantErr {
				t.Errorf("topUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockCoinProvider struct {
	err error
}

func (m mockCoinProvider) TopUpChan() chan types.TopUpRequest {
	tch := make(chan types.TopUpRequest)

	go func() {
		for tr := range tch {
			tr.ErrResp <- m.err
		}
	}()

	return tch
}

func (m mockCoinProvider) Stake(context.Context, string, string, *vega.Asset, *num.Uint, string) error {
	return m.err
}
