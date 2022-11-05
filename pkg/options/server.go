package options

import (
	"github.com/planetfall/gateway/pkg/config"
)

type ServerOptions struct {
	Env          string
	ServiceName  string
	ProjectID    string
	Port         string
	SvcConfigMap config.ServiceConfigMap
	Insecure     bool
}
