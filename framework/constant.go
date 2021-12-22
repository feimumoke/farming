package framework

const (
	NetProtocol    = "tcp"
	NetGrpcAddress = "0.0.0.0:28086"
	NetGwAddress   = "0.0.0.0:28088"
	NetGrpcPort    = "28086"
	NetGwPort      = "28088"
	LogPath        = "/var/log/farming/"
)

const (
	GRPC_GO_LOG_SEVERITY_LEVEL      = "GRPC_GO_LOG_SEVERITY_LEVEL"
	GRPC_GO_LOG_SEVERITY_LEVEL_VAL  = "INFO"
	GRPC_GO_LOG_VERBOSITY_LEVEL     = "GRPC_GO_LOG_VERBOSITY_LEVEL"
	GRPC_GO_LOG_VERBOSITY_LEVEL_VAL = "99"
)

const (
	SERVER_FRAMER_TARGET  = "dns:///server-farmer.default.svc.cluster.local"
	SERVER_FRAMER_PROJECT = "server-farmer"
	SERVER_FRAMER_SERVICE = "farmer"

	SERVER_GROUND_TARGET  = "dns:///server-ground.default.svc.cluster.local"
	SERVER_GROUND_PROJECT = "server-ground"
	SERVER_GROUND_SERVICE = "ground"
)
