package events

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	_ "github.com/aws/smithy-go/endpoints"
	"github.com/goccy/go-json"
	"github.com/saleh-ghazimoradi/Cartopher/config"
	"github.com/saleh-ghazimoradi/Cartopher/pkg/uploadProvider"
)

type WatermillEventPublisher struct {
	publisher message.Publisher
	queueName string
}

func (w *WatermillEventPublisher) Publish(eventType string, payload any, metadata map[string]string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), data)

	msg.Metadata.Set("event_type", eventType)

	for k, v := range metadata {
		msg.Metadata.Set(k, v)
	}

	return w.publisher.Publish(w.queueName, msg)
}

func (w *WatermillEventPublisher) Close() error {
	return w.publisher.Close()
}

func NewWatermillEventPublisher(ctx context.Context, cfg *config.Config) (*WatermillEventPublisher, error) {
	logger := watermill.NewStdLogger(false, false)

	awsConfig, err := uploadProvider.CreateAWSConfig(ctx, cfg.AWS.S3Endpoint, cfg.AWS.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create aws config: %w", err)
	}

	publisherConfig := sqs.PublisherConfig{
		AWSConfig: awsConfig,
		Marshaler: nil,
	}

	publisher, err := sqs.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create sqs publisher: %w", err)
	}

	return &WatermillEventPublisher{
		publisher: publisher,
		queueName: cfg.AWS.EventQueueName,
	}, nil
}
