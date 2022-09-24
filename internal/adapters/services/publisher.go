package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/eventual-agent/internal/pb"
	"github.com/josh-tracey/eventual-agent/internal/ports"
	"github.com/josh-tracey/notary"
	"github.com/josh-tracey/scribe"
	"google.golang.org/grpc"
)

var jwtTokenSecret = os.Getenv("JWT_TOKEN_SECRET")

type Publisher struct {
	logger         *scribe.Logger
	publishChannel chan *core.PeerEvent
	subs           *core.Adapter
}

func NewPublisher(subs ports.SubjectPort, logger *scribe.Logger, publishChannel chan *core.PeerEvent) (*Publisher, error) {
	value, err := subs.(*core.Adapter)
	if !err {
		return nil, errors.New("Invalid Subject Port")
	}
	return &Publisher{
			subs:           value,
			logger:         logger,
			publishChannel: publishChannel,
		},
		nil
}

func (p *Publisher) Publish(ctx context.Context, event *core.PeerEvent) error {

	conn, err := grpc.Dial(event.PeerServer, grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := pb.NewPublisherServiceClient(conn)
	defer conn.Close()

	c, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	token, err := notary.New(jwtTokenSecret).NewSignedToken()
	if err != nil {
		return err
	}

	_, err2 := client.Publish(
		c,
		&pb.EventPubRequest{
			Token: token,
			Data: &pb.CloudEvent{
				Id:      event.Event.ID,
				Type:    event.Event.Type,
				Subject: event.Event.Subject,
				Time:    event.Event.Time,
				Source:  event.Event.Source,
				Data:    &pb.CloudEvent_TextData{TextData: event.Event.Data},
			},
		}, grpc.FailFast(true))

	if err2 != nil {
		return err2
	}

	return nil
}

func (p *Publisher) Run() {

	for {
		select {
		case peerEvent := <-p.publishChannel:
			ok := p.subs.HasPeerId(peerEvent.PeerServer, false)
			if ok {
				start := time.Now()
				p.logger.Trace("Publishing Event => '%v' to Peer Server %v ", peerEvent.Event, peerEvent.PeerServer)
				err := p.Publish(context.Background(), peerEvent)
				if err != nil {
					p.logger.Error("Error publishing event: %v", err)
				}
				p.logger.Duration(start, "Publishing event to peer")
			}
		}
	}

}
