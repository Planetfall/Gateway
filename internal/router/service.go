package router

type ServicesConfig map[string]ServiceConfig

type ServiceConfig struct {
	Host     string
	Audience string

	QueueID        string
	LocationID     string
	SubscriptionID string

	WebsocketOrigins []string
}
