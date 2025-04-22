package console

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"goeasy.dev/events"
)

// Driver implements the events.Driver interface for console output
type Driver struct{}

// NewDriver creates a new console driver
func NewDriver() *Driver {
	return &Driver{}
}

// formatEvent formats an event for console output
func formatEvent(event events.Event) string {
	// Format the time in RFC3339
	timeStr := event.Time.Format(time.RFC3339)

	// Convert data to JSON string
	dataJSON, err := json.MarshalIndent(event.Data, "", "  ")
	if err != nil {
		dataJSON = []byte(fmt.Sprintf("error marshaling data: %v", err))
	}

	// Format metadata
	var metadataStr string
	if len(event.Metadata) > 0 {
		metadataJSON, err := json.MarshalIndent(event.Metadata, "", "  ")
		if err != nil {
			metadataStr = fmt.Sprintf("error marshaling metadata: %v", err)
		} else {
			metadataStr = string(metadataJSON)
		}
	}

	// Build the output string
	output := fmt.Sprintf("Event:\n"+
		"  ID:      %s\n"+
		"  Source:  %s\n"+
		"  Type:    %s\n"+
		"  Time:    %s\n"+
		"  Data:    %s",
		event.ID,
		event.Source,
		event.Type,
		timeStr,
		string(dataJSON))

	if len(event.Metadata) > 0 {
		output += fmt.Sprintf("\n  Metadata: %s", metadataStr)
	}

	return output
}

// Publish prints a single event to the console
func (d *Driver) Publish(ctx context.Context, event events.Event) error {
	fmt.Println(formatEvent(event))
	return nil
}

// PublishBatch prints multiple events to the console
func (d *Driver) PublishBatch(ctx context.Context, events []events.Event) error {
	if len(events) == 0 {
		return nil
	}

	fmt.Printf("Publishing batch of %d events:\n\n", len(events))
	for i, event := range events {
		fmt.Printf("Event %d/%d:\n%s\n\n", i+1, len(events), formatEvent(event))
	}
	return nil
}

// Close implements the Driver interface
func (d *Driver) Close() error {
	return nil
}
