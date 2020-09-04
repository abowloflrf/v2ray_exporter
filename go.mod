module github.com/abowloflrf/v2ray-exporter

go 1.13

require (
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/cobra v0.0.7
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20200831180312-196b9ba8737a // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.31.1
	v2ray.com/core v4.19.1+incompatible
)

replace v2ray.com/core => github.com/v2ray/v2ray-core v4.27.5+incompatible
