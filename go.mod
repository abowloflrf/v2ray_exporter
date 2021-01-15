module github.com/abowloflrf/v2ray-exporter

go 1.13

require (
	github.com/prometheus/client_golang v1.9.0
	github.com/spf13/cobra v0.0.7
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.35.0
	v2ray.com/core v4.19.1+incompatible
)

replace v2ray.com/core => github.com/v2ray/v2ray-core v4.27.5+incompatible
