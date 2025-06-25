package client

import (
	"context"

	"github.com/korolev-n/gExchange/shared/api"
	"google.golang.org/grpc"
)

type ExchangerClient struct {
	client api.ExchangerServiceClient
}

func NewExchangerClient(conn *grpc.ClientConn) *ExchangerClient {
	return &ExchangerClient{
		client: api.NewExchangerServiceClient(conn),
	}
}

func (e *ExchangerClient) GetRates(ctx context.Context) (map[string]float64, error) {
	res, err := e.client.GetRates(ctx, &api.Empty{})
	if err != nil {
		return nil, err
	}
	return res.Rates, nil
}
