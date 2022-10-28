package server

import (
	musicResearcherPb "github.com/Dadard29/planetfall/gateway/pkg/musicresearcher"
)

type clients struct {
	musicResearcher musicResearcherPb.MusicResearcherClient
}

func newClients(conns connections) *clients {
	return &clients{
		musicResearcher: musicResearcherPb.NewMusicResearcherClient(
			conns[conn_MusicResearcher].grpcConn),
	}
}
