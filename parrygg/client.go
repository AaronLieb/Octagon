package parrygg

import (
	"context"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/AaronLieb/octagon/parrygg/pb"
)

type Client struct {
	conn                *grpc.ClientConn
	apiKey              string
	TournamentService   pb.TournamentServiceClient
	UserService         pb.UserServiceClient
	EventService        pb.EventServiceClient
	EntrantService      pb.EntrantServiceClient
	BracketService      pb.BracketServiceClient
	PhaseService        pb.PhaseServiceClient
	MatchService        pb.MatchServiceClient
	MatchGameService    pb.MatchGameServiceClient
	GameService         pb.GameServiceClient
	NotificationService pb.NotificationServiceClient
}

func loggingInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		log.Infof("gRPC Request: %s", method)
		log.Infof("Metadata: %+v", md)
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

func NewClient(apiKey string) (*Client, error) {
	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.NewClient(
		"api.parry.gg:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(loggingInterceptor),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:                conn,
		apiKey:              apiKey,
		TournamentService:   pb.NewTournamentServiceClient(conn),
		UserService:         pb.NewUserServiceClient(conn),
		EventService:        pb.NewEventServiceClient(conn),
		EntrantService:      pb.NewEntrantServiceClient(conn),
		BracketService:      pb.NewBracketServiceClient(conn),
		PhaseService:        pb.NewPhaseServiceClient(conn),
		MatchService:        pb.NewMatchServiceClient(conn),
		MatchGameService:    pb.NewMatchGameServiceClient(conn),
		GameService:         pb.NewGameServiceClient(conn),
		NotificationService: pb.NewNotificationServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) WithAuth(ctx context.Context) context.Context {
	md := metadata.Pairs("X-API-KEY", c.apiKey)
	return metadata.NewOutgoingContext(ctx, md)
}
