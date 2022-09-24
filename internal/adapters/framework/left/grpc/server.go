package grpc

import (
	"context"
	"errors"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/eventual-agent/internal/pb"
	"github.com/josh-tracey/eventual-agent/internal/ports"
	"github.com/josh-tracey/notary"
	"github.com/josh-tracey/scribe"
)

var jwtTokenSecret = os.Getenv("JWT_TOKEN_SECRET")

type Adapter struct {
	pb.UnimplementedClientServiceServer
	core           *core.Adapter
	logger         *scribe.Logger
	publishChannel chan *core.PeerEvent
	subsChannel    chan *core.PeerRequest
}

func New(
	c ports.SubjectPort,
	logger *scribe.Logger,
	publishChannel chan *core.PeerEvent,
	subsChannel chan *core.PeerRequest,
) *Adapter {

	value, ok := c.(*core.Adapter)
	if !ok {
		panic("Invalid core adapter")
	}
	return &Adapter{
		logger:         logger,
		core:           value,
		publishChannel: publishChannel,
		subsChannel:    subsChannel,
	}
}

func (a *Adapter) Subscribe(ctx context.Context, req *pb.EventSubRequest) (*pb.EventSubResponse, error) {
	valid, err := notary.New(jwtTokenSecret).VerifyToken(req.Token)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("Invalid token")
	}

	a.logger.Debug("PeerServer: %s", req.PeerServer)

	a.subsChannel <- &core.PeerRequest{
		PeerAddr:  req.PeerServer,
		Channel:   req.Channel,
		Ephemeral: false,
	}

	a.core.AddPeer(req.PeerServer, req.Channel, false)

	return &pb.EventSubResponse{SubscriptionId: ""}, nil

}

func (a *Adapter) Publish(ctx context.Context, req *pb.EventPubRequest) (*pb.EventPubResponse, error) {
	valid, err := notary.New(jwtTokenSecret).VerifyToken(req.Token)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("Invalid token")
	}

	event := core.CloudEvent{
		ID:      req.Data.Id,
		Source:  req.Data.Source,
		Type:    req.Data.Type,
		Subject: req.Data.Subject,
		Data:    req.Data.GetTextData(),
		Time:    req.Data.Time,
	}

	go func() {
		for peerServer := range a.core.GetPeerServers() {
			a.publishChannel <- &core.PeerEvent{
				PeerServer: peerServer,
				Event:      event,
			}
		}
	}()

	return &pb.EventPubResponse{SubscriptionId: req.SubscriptionId}, nil
}

func (a *Adapter) Run() error {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	a.logger.Info("gRPC Server Listening on 0.0.0.0:9090")
	pb.RegisterClientServiceServer(s, a)
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
