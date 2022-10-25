package server

import (
	"log"

	"github.com/spf13/viper"
)

type connections struct {
	musicResearcher *connection
}

func newConnections() (*connections, error) {
	conns := &connections{}

	// music-researcher
	musicResearcher, err := newConnection(
		viper.GetString("services.music-researcher.host"),
		viper.GetString("services.music-researcher.audience"),
	)
	if err != nil {
		log.Printf("failed setting up connection to music-researcher: %v\n", err)
		return nil, err
	}

	conns.musicResearcher = musicResearcher

	return conns, nil
}
