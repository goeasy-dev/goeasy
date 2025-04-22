package events

import (
	"context"
	"sync"
	"time"
)

// Event represents a generic event that can be published
type Event struct {
	// ID is a unique identifier for the event
	ID string
	// Source is the origin of the event (e.g., "com.myapp.users")
	Source string
	// Type is the type of event (e.g., "user.created")
	Type string
	// Time is when the event occurred
	Time time.Time
	// Data is the event payload
	Data interface{}
	// Metadata contains additional event information
	Metadata map[string]string
}

// Driver defines the interface that event publishing drivers must implement
type Driver interface {
	// Publish sends an event to the event bus
	Publish(ctx context.Context, event Event) error
	// PublishBatch sends multiple events to the event bus
	PublishBatch(ctx context.Context, events []Event) error
	// Close closes the driver and releases any resources
	Close() error
}

// Registry manages event publishing drivers
type Registry struct {
	drivers map[string]Driver
	mu      sync.RWMutex
}

// globalRegistry is the default registry used by the package
var globalRegistry = NewRegistry()

// NewRegistry creates a new driver registry
func NewRegistry() *Registry {
	return &Registry{
		drivers: make(map[string]Driver),
	}
}

// RegisterDriver registers a new driver with the registry
func (r *Registry) RegisterDriver(name string, driver Driver) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.drivers[name] = driver
}

// GetDriver returns a driver by name
func (r *Registry) GetDriver(name string) (Driver, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	driver, ok := r.drivers[name]
	return driver, ok
}

// Close closes all registered drivers
func (r *Registry) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var lastErr error
	for _, driver := range r.drivers {
		if err := driver.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// RegisterDriver registers a driver with the global registry
func RegisterDriver(name string, driver Driver) {
	globalRegistry.RegisterDriver(name, driver)
}

// Close closes the global registry
func Close() error {
	return globalRegistry.Close()
}

// Publish publishes a single event using the default driver
func Publish(ctx context.Context, event Event) error {
	return PublishWithDriver(ctx, "default", event)
}

// PublishWithDriver publishes a single event using the specified driver
func PublishWithDriver(ctx context.Context, driverName string, event Event) error {
	driver, ok := globalRegistry.GetDriver(driverName)
	if !ok {
		return ErrDriverNotFound
	}

	return driver.Publish(ctx, event)
}

// PublishBatch publishes multiple events using the default driver
func PublishBatch(ctx context.Context, events []Event) error {
	return PublishBatchWithDriver(ctx, "default", events)
}

// PublishBatchWithDriver publishes multiple events using the specified driver
func PublishBatchWithDriver(ctx context.Context, driverName string, events []Event) error {
	driver, ok := globalRegistry.GetDriver(driverName)
	if !ok {
		return ErrDriverNotFound
	}

	return driver.PublishBatch(ctx, events)
}

// SetDefaultDriver sets the default driver to use
func SetDefaultDriver(name string) error {
	_, ok := globalRegistry.GetDriver(name)
	if !ok {
		return ErrDriverNotFound
	}
	return nil
}
