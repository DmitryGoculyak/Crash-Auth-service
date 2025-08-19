package currency

import (
	"context"
	"fmt"
	"log"
	"time"

	proto "Crash-Auth-service/pkg/proto/currency"
	"google.golang.org/grpc"
)

type CurrencyClient struct {
	conn   *grpc.ClientConn
	client proto.CurrencyServiceClient
}

func CurrencyAdapter(cfg *CurrencyConfig) (*CurrencyClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	maxMsgSize := 10500000

	dialOption := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		),
	}

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	conn, err := grpc.DialContext(ctx, address, dialOption...)
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}

	client := proto.NewCurrencyServiceClient(conn)

	return &CurrencyClient{
		conn,
		client,
	}, nil
}

func (s *CurrencyClient) CheckCurrencyExists(ctx context.Context, currencyCode string) error {
	_, err := s.client.GetCurrencies(ctx, &proto.GetCurrenciesRequest{CurrencyCode: currencyCode})
	if err != nil {
		return err
	}
	return nil
}
