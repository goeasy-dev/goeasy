package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"goeasy.dev/errors"
	"goeasy.dev/events"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

// Config holds the configuration for the AWS EventBridge driver
type Config struct {
	// BusName is the name of the event bus (optional)
	BusName string
}

// Driver implements the events.Driver interface for AWS EventBridge
type Driver struct {
	config Config
	client *eventbridge.Client
}

// NewDriver creates a new AWS EventBridge driver
func NewDriver(ctx context.Context, cfg Config) (*Driver, error) {
	// Load the AWS SDK configuration from the environment and shared config
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load AWS configuration")
	}

	return &Driver{
		client: eventbridge.NewFromConfig(awsCfg),
		config: cfg,
	}, nil
}

// Close closes the driver
func (d *Driver) Close() error {
	// Nothing to close for EventBridge
	return nil
}

// Publish sends an event to EventBridge
func (d *Driver) Publish(ctx context.Context, event events.Event) error {
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
				EventBusName: aws.String(d.config.BusName),
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
	result, err := d.client.PutEvents(ctx, ebEvent)
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
func (d *Driver) PublishBatch(ctx context.Context, events []events.Event) error {
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
			EventBusName: aws.String(d.config.BusName),
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
	result, err := d.client.PutEvents(ctx, &eventbridge.PutEventsInput{
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
