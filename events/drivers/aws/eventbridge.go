package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"goeasy.dev/errors"
	"goeasy.dev/events"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

// Config holds the configuration for the AWS EventBridge driver
type Config struct {
	// Region is the AWS region to use
	Region string
	// BusName is the name of the event bus (optional)
	BusName string
}

// Driver implements the events.Driver interface for AWS EventBridge
type Driver struct {
	config Config
	client *eventbridge.Client
}

// NewDriver creates a new AWS EventBridge driver
func NewDriver() *Driver {
	return &Driver{}
}

// Initialize sets up the AWS EventBridge client
func (d *Driver) Initialize(ctx context.Context, config interface{}) error {
	cfg, ok := config.(Config)
	if !ok {
		return errors.New("invalid config type for AWS EventBridge driver")
	}
	d.config = cfg

	// Create AWS config
	awsCfg := aws.Config{
		Region: cfg.Region,
	}

	// Create EventBridge client
	d.client = eventbridge.NewFromConfig(awsCfg)
	return nil
}

// CreatePublisher creates a new EventBridge publisher
func (d *Driver) CreatePublisher(ctx context.Context, config interface{}) (events.Publisher, error) {
	if d.client == nil {
		return nil, errors.New("driver not initialized")
	}
	return &Publisher{
		client: d.client,
		config: d.config,
	}, nil
}

// Close closes the driver
func (d *Driver) Close() error {
	// Nothing to close for EventBridge
	return nil
}

// Publisher implements the events.Publisher interface for AWS EventBridge
type Publisher struct {
	client *eventbridge.Client
	config Config
}

// Publish sends an event to EventBridge
func (p *Publisher) Publish(ctx context.Context, event events.Event) error {
	// Convert event data to JSON
	data, err := json.Marshal(event.Data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal event data")
	}

	// Create EventBridge event
	ebEvent := &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Source:       aws.String(event.Source),
				DetailType:   aws.String(event.Type),
				Detail:       aws.String(string(data)),
				EventBusName: aws.String(p.config.BusName),
				Time:         aws.Time(event.Time),
			},
		},
	}

	// Add metadata as additional attributes if present
	if len(event.Metadata) > 0 {
		ebEvent.Entries[0].Resources = make([]string, 0, len(event.Metadata))
		for k, v := range event.Metadata {
			ebEvent.Entries[0].Resources = append(ebEvent.Entries[0].Resources, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// Send event
	result, err := p.client.PutEvents(ctx, ebEvent)
	if err != nil {
		return errors.Wrap(err, "failed to publish event to EventBridge")
	}

	// Check for failed entries
	if result.FailedEntryCount > 0 {
		return errors.New(fmt.Sprintf("failed to publish %d events", result.FailedEntryCount))
	}

	return nil
}

// PublishBatch sends multiple events to EventBridge
func (p *Publisher) PublishBatch(ctx context.Context, events []events.Event) error {
	if len(events) == 0 {
		return nil
	}

	// Convert events to EventBridge format
	ebEvents := make([]types.PutEventsRequestEntry, len(events))
	for i, event := range events {
		// Convert event data to JSON
		data, err := json.Marshal(event.Data)
		if err != nil {
			return errors.Wrap(err, "failed to marshal event data")
		}

		ebEvents[i] = types.PutEventsRequestEntry{
			Source:       aws.String(event.Source),
			DetailType:   aws.String(event.Type),
			Detail:       aws.String(string(data)),
			EventBusName: aws.String(p.config.BusName),
			Time:         aws.Time(event.Time),
		}

		// Add metadata as additional attributes if present
		if len(event.Metadata) > 0 {
			ebEvents[i].Resources = make([]string, 0, len(event.Metadata))
			for k, v := range event.Metadata {
				ebEvents[i].Resources = append(ebEvents[i].Resources, fmt.Sprintf("%s=%s", k, v))
			}
		}
	}

	// Send events
	result, err := p.client.PutEvents(ctx, &eventbridge.PutEventsInput{
		Entries: ebEvents,
	})
	if err != nil {
		return errors.Wrap(err, "failed to publish events to EventBridge")
	}

	// Check for failed entries
	if result.FailedEntryCount > 0 {
		return errors.New(fmt.Sprintf("failed to publish %d events", result.FailedEntryCount))
	}

	return nil
}

// Close closes the publisher
func (p *Publisher) Close() error {
	// Nothing to close for EventBridge
	return nil
}
