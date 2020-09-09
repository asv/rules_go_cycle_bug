package pubsub

import (
	"context"

	"github.com/golang/protobuf/proto"
)

type Message proto.Message
type MessageHandler func(ctx context.Context, msgData []byte)

type PublishOption interface{}

type Publisher interface {
	Publish(ctx context.Context, subject string, msg proto.Message, pubOpts ...PublishOption) error
}

type SubscriptionOption interface{}

type Subscriber interface {
	Subscribe(ctx context.Context, subject string, handler MessageHandler, subOpts ...SubscriptionOption) error
}

type PubSuber interface {
	Publisher
	Subscriber
}
