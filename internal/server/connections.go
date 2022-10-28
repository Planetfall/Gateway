package server

import (
	"log"
)

type ConnectionConfigList map[string]ConnectionConfig
type ConnectionConfig struct {
	Host     string `mapstructure:"host"`
	Audience string `mapstructure:"audience"`
}

const (
	conn_MusicResearcher = "music-researcher"
)

type connections map[string]*connection

func newConnections(
	connCfgList ConnectionConfigList, insecure bool) (connections, error) {

	conns := connections{}

	for connectionName, connectionConfig := range connCfgList {
		var err error
		var conn *connection
		if insecure == false {
			conn, err = newConnection(
				connectionConfig.Host, connectionConfig.Audience)
		}

		if insecure == true {
			conn, err = newConnectionInsecure(connectionConfig.Host)
		}

		if err != nil {
			log.Printf(
				"failed to set up connection to %s: %v\n", connectionName, err)
		}
		conns[connectionName] = conn
	}

	return conns, nil
}

func (cs connections) Close() error {
	for _, conn := range cs {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
