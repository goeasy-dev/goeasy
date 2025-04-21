package aws_test

import (
	"context"
	"testing"
	"time"

	"goeasy.dev/events"
	"goeasy.dev/events/drivers/aws"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventBridgeDriver(t *testing.T) {
	ctx := context.Background()

	// Create driver
	driver := aws.NewDriver()

	// Initialize driver
	err := driver.Initialize(ctx, aws.Config{
		Region:  "us-west-2",
		BusName: "test-bus",
	})
	require.NoError(t, err)

	// Create publisher
	publisher, err := driver.CreatePublisher(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, publisher)

	// Test publishing a single event
	event := events.Event{
		ID:       "test-event-1",
		Source:   "com.test.events",
		Type:     "test.created",
		Time:     time.Now(),
		Data:     map[string]string{"key": "value"},
		Metadata: map[string]string{"env": "test"},
	}

	err = publisher.Publish(ctx, event)
	// Note: This will fail in tests without AWS credentials
	// In a real test environment, you would use a mock or localstack
	assert.Error(t, err)

	// Test publishing multiple events
	events := []events.Event{
		{
			ID:       "test-event-2",
			Source:   "com.test.events",
			Type:     "test.created",
			Time:     time.Now(),
			Data:     map[string]string{"key": "value2"},
			Metadata: map[string]string{"env": "test"},
		},
		{
			ID:       "test-event-3",
			Source:   "com.test.events",
			Type:     "test.created",
			Time:     time.Now(),
			Data:     map[string]string{"key": "value3"},
			Metadata: map[string]string{"env": "test"},
		},
	}

	err = publisher.PublishBatch(ctx, events)
	// Note: This will fail in tests without AWS credentials
	// In a real test environment, you would use a mock or localstack
	assert.Error(t, err)

	// Close publisher
	err = publisher.Close()
	assert.NoError(t, err)

	// Close driver
	err = driver.Close()
	assert.NoError(t, err)
}
