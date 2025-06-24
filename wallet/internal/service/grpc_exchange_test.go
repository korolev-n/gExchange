package service_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/korolev-n/gExchange/exchanger/api"
	"github.com/korolev-n/gExchange/wallet/internal/cache"
	"github.com/korolev-n/gExchange/wallet/internal/domain"
	"github.com/korolev-n/gExchange/wallet/internal/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type mockExchangerServer struct {
	api.UnimplementedExchangerServiceServer
}

type mockWalletRepo struct {
	GetBalancesFunc   func(ctx context.Context, userID int64) ([]domain.Wallet, error)
	UpdateBalanceFunc func(ctx context.Context, userID int64, currency string, delta float64) error
}

func (m *mockWalletRepo) GetBalances(ctx context.Context, userID int64) ([]domain.Wallet, error) {
	return m.GetBalancesFunc(ctx, userID)
}

func (m *mockWalletRepo) UpdateBalance(ctx context.Context, userID int64, currency string, delta float64) error {
	return m.UpdateBalanceFunc(ctx, userID, currency, delta)
}

// Мок-сервера
func (m *mockExchangerServer) GetRates(ctx context.Context, _ *api.Empty) (*api.RatesResponse, error) {
	return &api.RatesResponse{
		Rates: map[string]float64{
			"USD": 1.0,
			"EUR": 1.1765,
			"RUB": 1.0,
		},
	}, nil
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func dialer() func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func startBufconnServer(t *testing.T) *grpc.ClientConn {
	lis := bufconn.Listen(bufSize)
	t.Cleanup(func() { lis.Close() })

	s := grpc.NewServer()
	api.RegisterExchangerServiceServer(s, &mockExchangerServer{})

	go func() {
		if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("Server exited with error: %v", err)
		}
	}()
	t.Cleanup(func() { s.Stop() })

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Современная замена WithInsecure
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Logf("Failed to close connection: %v", err)
		}
	})

	return conn
}

type grpcExchangerAdapter struct {
	client api.ExchangerServiceClient
}

func (g *grpcExchangerAdapter) GetRates(ctx context.Context) (map[string]float64, error) {
	resp, err := g.client.GetRates(ctx, &api.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Rates, nil
}

func TestWalletService_Exchange_Integration(t *testing.T) {
	conn := startBufconnServer(t)
	defer conn.Close()

	grpcClient := api.NewExchangerServiceClient(conn)
	adapter := &grpcExchangerAdapter{client: grpcClient}

	mockRepo := &mockWalletRepo{
		GetBalancesFunc: func(ctx context.Context, userID int64) ([]domain.Wallet, error) {
			return []domain.Wallet{
				{UserID: userID, Currency: "USD", Balance: 100},
				{UserID: userID, Currency: "EUR", Balance: 0},
			}, nil
		},
		UpdateBalanceFunc: func(ctx context.Context, userID int64, currency string, delta float64) error {
			return nil
		},
	}

	cache := cache.NewExchangeRateCache(30 * time.Second)
	service := service.NewWalletService(mockRepo, adapter, cache)

	bal, exchanged, err := service.Exchange(context.Background(), 1, "USD", "EUR", 100)
	require.NoError(t, err)
	require.InDelta(t, 85.0, exchanged, 0.01)
	require.InDelta(t, 0.0, bal["USD"], 0.01)
	require.InDelta(t, 85.0, bal["EUR"], 0.01)
}
