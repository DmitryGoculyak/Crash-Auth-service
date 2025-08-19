package billing

import (
	"context"
	"fmt"
	"log"
	"time"

	proto "Crash-Auth-service/pkg/proto/billing"
	"google.golang.org/grpc"
)

type BillingClient struct {
	conn   *grpc.ClientConn
	client proto.BillingServiceClient
}

func BillingAdapter(cfg *BillingConfig) (*BillingClient, error) {
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
		log.Printf("Did not connect: %v", err)
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}

	client := proto.NewBillingServiceClient(conn)

	return &BillingClient{
		conn,
		client,
	}, nil
}

func (s *BillingClient) CreateWallet(ctx context.Context, userID, currencyCode string) error {
	_, err := s.client.CreateWallet(ctx, &proto.CreateWalletRequest{
		UserId:       userID,
		CurrencyCode: currencyCode,
	})
	if err != nil {
		return err
	}
	return nil
}
