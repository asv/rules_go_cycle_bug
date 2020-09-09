package pubsub

import (
	"context"
	"fmt"

	"github.com/nats-io/stan.go"

	"github.com/golang/protobuf/proto"
)

type StanPubsub struct {
	stanConn stan.Conn

	clusterID string
	clientID  string
}

var (
	StanConnectionURL = stan.NatsURL
)

func NewStanPubsub(clusterID, clientID string, stanOpts ...stan.Option) (*StanPubsub, error) {
	stanConn, err := stan.Connect(clusterID, clientID, stanOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create STAN connection: %w", err)
	}

	p := &StanPubsub{
		stanConn: stanConn,

		clientID:  clientID,
		clusterID: clusterID,
	}

	return p, nil
}

func (p *StanPubsub) Close() error {
	return p.stanConn.Close()
}

func (p *StanPubsub) Publish(ctx context.Context, subject string, msg proto.Message, pubOpts ...PublishOption) error {
	b, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := p.stanConn.Publish(subject, b); err != nil {
		return fmt.Errorf("failed to publish message via STAN: %w", err)
	}

	return nil
}

func (p *StanPubsub) Subscribe(ctx context.Context, subject string, handler MessageHandler, subOpts ...SubscriptionOption) error {
	durableName := fmt.Sprintf("%s-%s", p.clientID, subject)

	fn := func(msg *stan.Msg) {
		handler(context.Background(), msg.MsgProto.Data)
	}

	stanSubOpts := []stan.SubscriptionOption{stan.DurableName(durableName)}
	for _, opt := range subOpts {
		stanSubOpt := opt.(stan.SubscriptionOption)
		stanSubOpts = append(stanSubOpts, stanSubOpt)
	}

	if _, err := p.stanConn.Subscribe(subject, fn, stanSubOpts...); err != nil {
		return fmt.Errorf("failed to subscribe (subject=%q) via NATS: %w", subject, err)
	}

	return nil
}
