package main

import (
	"context"
	"log"
	"time"

	"goeasy.dev/events"
	"goeasy.dev/events/drivers/aws"
)

func main() {
	ctx := context.Background()

	// Register the AWS EventBridge driver
	driver := aws.NewDriver()
	events.RegisterDriver("eventbridge", driver)

	// Initialize the driver
	err := events.InitializeDriver(ctx, "eventbridge", aws.Config{
		Region:  "us-west-2",
		BusName: "my-event-bus",
	})
	if err != nil {
		log.Fatalf("Failed to initialize driver: %v", err)
	}

	// Set as default driver
	err = events.SetDefaultDriver("eventbridge")
	if err != nil {
		log.Fatalf("Failed to set default driver: %v", err)
	}

	// Create an event
	event := events.Event{
		ID:     "event-1",
		Source: "com.myapp.users",
		Type:   "user.created",
		Time:   time.Now(),
		Data: map[string]interface{}{
			"userId":   "123",
			"username": "johndoe",
			"email":    "john@example.com",
		},
		Metadata: map[string]string{
			"environment": "production",
			"version":     "1.0",
		},
	}

	// Publish the event using the default driver
	err = events.Publish(ctx, event)
	if err != nil {
		log.Fatalf("Failed to publish event: %v", err)
	}

	// Create multiple events
	eventBatch := []events.Event{
		{
			ID:     "event-2",
			Source: "com.myapp.orders",
			Type:   "order.created",
			Time:   time.Now(),
			Data: map[string]interface{}{
				"orderId": "456",
				"amount":  99.99,
			},
			Metadata: map[string]string{
				"environment": "production",
				"version":     "1.0",
			},
		},
		{
			ID:     "event-3",
			Source: "com.myapp.orders",
			Type:   "order.shipped",
			Time:   time.Now(),
			Data: map[string]interface{}{
				"orderId":     "456",
				"trackingNum": "TRACK123",
			},
			Metadata: map[string]string{
				"environment": "production",
				"version":     "1.0",
			},
		},
	}

	// Publish multiple events using the default driver
	err = events.PublishBatch(ctx, eventBatch)
	if err != nil {
		log.Fatalf("Failed to publish events: %v", err)
	}

	// Example of using a specific driver
	err = events.PublishWithDriver(ctx, "eventbridge", event)
	if err != nil {
		log.Fatalf("Failed to publish event with specific driver: %v", err)
	}

	log.Println("Successfully published events")

	// Clean up
	err = events.Close()
	if err != nil {
		log.Printf("Error closing events: %v", err)
	}
}
