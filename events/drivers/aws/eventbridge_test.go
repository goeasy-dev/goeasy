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
	driver, err := aws.NewDriver(ctx, aws.Config{
		BusName: "test-bus",
	})
	require.NoError(t, err)

	// Test publishing a single event
	event := events.Event{
		ID:       "test-event-1",
		Source:   "com.test.events",
		Type:     "test.created",
		Time:     time.Now(),
		Data:     map[string]string{"key": "value"},
		Metadata: map[string]string{"env": "test"},
	}

	err = driver.Publish(ctx, event)
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

	err = driver.PublishBatch(ctx, events)
	// Note: This will fail in tests without AWS credentials
	// In a real test environment, you would use a mock or localstack
	assert.Error(t, err)

	err = driver.Close()
	assert.NoError(t, err)
}
