package config

type ServiceConfigMap map[string]ServiceConfig

type ServiceConfig struct {
	// grpc
	Host     string `mapstructure:"host"`
	Audience string `mapstructure:"audience"`

	// task queues
	QueueID    string `mapstructure:"queue_id"`
	LocationID string `mapstructure:"location_id"`

	// pubsub
	SubscriptionID string `mapstructure:"subscription_id"`

	// websockets
	WebsocketOrigins []string `mapstructure:"websocket_origins"`
}
