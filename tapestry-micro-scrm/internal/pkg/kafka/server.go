package pkgKafka

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/servers"
)

type ServerGroup struct {
	servers []servers.Server
}

func (s *ServerGroup) Register(server servers.Server) {
	s.servers = append(s.servers, server)
}

func (s *ServerGroup) Servers() []servers.Server {
	return s.servers
}
