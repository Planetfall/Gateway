package config

type ServiceConfigMap map[string]ServiceConfig

type ServiceConfig struct {
	Host     string `mapstructure:"host"`
	Audience string `mapstructure:"audience"`

	QueueID    string `mapstructure:"queue_id"`
	LocationID string `mapstructure:"location_id"`
}
