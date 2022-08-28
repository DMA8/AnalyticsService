package grpc

import (
	"gitlab.com/g6834/team31/analytics/internal/domain/models"
	"gitlab.com/g6834/team31/auth/pkg/grpc_auth"
	"gitlab.com/g6834/team31/auth/pkg/logging"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client grpc_auth.AuthClient
	conn   *grpc.ClientConn
	l      *logging.Logger
}

func New(ctx context.Context, host, port string, l *logging.Logger) (*Client, error) {
	connStr := host + port
	conn, err := grpc.DialContext(ctx, connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := grpc_auth.NewAuthClient(conn)
	return &Client{
		client: client,
		conn:   conn,
		l: l,
	}, nil
}

func (c *Client) Stop() error {
	return c.conn.Close()
}

func (c *Client) Validate(ctx context.Context, in models.JWTTokens) (models.ValidateResponse, error) {
	response, err := c.client.Validate(ctx, &grpc_auth.Credential{
		AccessToken:  in.Access,
		RefreshToken: in.Refresh,
	},
	)
	if err != nil {
		c.l.Debug().Msgf("grpcClient.Validate couldn't validate jwt tokens")
		return models.ValidateResponse{}, err
	}
	c.l.Info().Msg("auth succeed")
	return models.ValidateResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		Login:        response.Login,
		Success:      response.Success,
		IsUpdate:     response.IsUpdate,
	}, nil
}
