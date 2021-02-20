package main

import (
	"errors"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	v2stats "github.com/v2fly/v2ray-core/v4/app/stats/command"
)

const (
	namespace = "v2ray"
)

type Exporter struct {
	inboundTagTrafficUplink   *prometheus.Desc
	inboundTagTrafficDownlink *prometheus.Desc
	userTrafficUplink         *prometheus.Desc
	userTrafficDownlink       *prometheus.Desc

	numGoroutine *prometheus.Desc
	numGC        *prometheus.Desc
	alloc        *prometheus.Desc
	totalAlloc   *prometheus.Desc
	sys          *prometheus.Desc
	mallocs      *prometheus.Desc
	frees        *prometheus.Desc
	liveObjects  *prometheus.Desc
	pauseTotalNs *prometheus.Desc
	uptime       *prometheus.Desc
}

type SingleF64Stat struct {
	Name  string
	Value float64
	Type  prometheus.ValueType
	Key   string
	Tag   string
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.inboundTagTrafficUplink
	ch <- e.inboundTagTrafficDownlink
	ch <- e.userTrafficUplink
	ch <- e.userTrafficDownlink
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	usageStats, err := v2c.QueryStats("")
	if err != nil {
		logger.Warnln("get usage stats from v2ray", "error", err)
	} else {
		if err := e.parseUsageStats(ch, usageStats); err != nil {
			logger.Warnln("parse usage stats from v2ray", "error", err)
		}
	}

	sysStats, err := v2c.GetSysStats()
	if err != nil {
		logger.Warnln("get sys sysStats from v2ray", "error", err)
		return
	}
	logger.Debugln("collected sys stats", "data", sysStats)
	ch <- prometheus.MustNewConstMetric(e.numGoroutine, prometheus.GaugeValue, float64(sysStats.NumGoroutine))
	ch <- prometheus.MustNewConstMetric(e.numGC, prometheus.CounterValue, float64(sysStats.NumGC))
	ch <- prometheus.MustNewConstMetric(e.alloc, prometheus.GaugeValue, float64(sysStats.Alloc))
	ch <- prometheus.MustNewConstMetric(e.totalAlloc, prometheus.GaugeValue, float64(sysStats.TotalAlloc))
	ch <- prometheus.MustNewConstMetric(e.sys, prometheus.GaugeValue, float64(sysStats.Sys))
	ch <- prometheus.MustNewConstMetric(e.mallocs, prometheus.GaugeValue, float64(sysStats.Mallocs))
	ch <- prometheus.MustNewConstMetric(e.frees, prometheus.GaugeValue, float64(sysStats.Frees))
	ch <- prometheus.MustNewConstMetric(e.liveObjects, prometheus.GaugeValue, float64(sysStats.LiveObjects))
	ch <- prometheus.MustNewConstMetric(e.pauseTotalNs, prometheus.CounterValue, float64(sysStats.PauseTotalNs))
	ch <- prometheus.MustNewConstMetric(e.uptime, prometheus.CounterValue, float64(sysStats.Uptime))

}

func (e *Exporter) parseStatsItem(s *v2stats.Stat) (*SingleF64Stat, error) {
	// "user>>>i@ruofeng.me>>>traffic>>>downlink"
	// "inbound>>>api>>>traffic>>>uplink"
	// "inbound>>>vmess-ws-in>>>traffic>>>downlink"
	d := strings.Split(s.Name, ">>>")
	if len(d) != 4 {
		return nil, errors.New("invalid stats [length]")
	}

	tag := d[1]
	link := d[3]

	if d[0] == "inbound" {
		return &SingleF64Stat{
			Name:  "inbound_tag_traffic_" + link,
			Value: float64(s.Value),
			Type:  prometheus.CounterValue,
			Key:   "tag",
			Tag:   tag,
		}, nil
	}

	if d[0] == "user" {
		return &SingleF64Stat{
			Name:  "user_traffic_" + link,
			Value: float64(s.Value),
			Type:  prometheus.CounterValue,
			Key:   "user",
			Tag:   tag,
		}, nil
	}
	return nil, errors.New("invalid stats [type]")
}

func (e *Exporter) parseUsageStats(ch chan<- prometheus.Metric, stats []*v2stats.Stat) error {
	itemsMetrics := map[string]*prometheus.Desc{
		"inbound_tag_traffic_uplink":   e.inboundTagTrafficUplink,
		"inbound_tag_traffic_downlink": e.inboundTagTrafficDownlink,
		"user_traffic_uplink":          e.userTrafficUplink,
		"user_traffic_downlink":        e.userTrafficDownlink,
	}
	for m, d := range itemsMetrics {
		for _, t := range stats {
			single, err := e.parseStatsItem(t)
			if err != nil {
				logger.Warnln("parse usage stats", "error", err)
				break
			}
			if single.Name == m {
				ch <- prometheus.MustNewConstMetric(d, single.Type, single.Value, single.Tag)
				logger.Debugln("collected usage stats", "data", single)
			}
		}
	}
	return nil
}

func NewExporter() *Exporter {
	return &Exporter{
		inboundTagTrafficUplink: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "inbound_tag_traffic_uplink"),
			"System inbound uplink traffic group by tag.",
			[]string{"tag"}, nil,
		),
		inboundTagTrafficDownlink: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "inbound_tag_traffic_downlink"),
			"System inbound downlink traffic group by tag.",
			[]string{"tag"}, nil,
		),
		userTrafficUplink: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "user_traffic_uplink"),
			"User uplink traffic.",
			[]string{"user"}, nil,
		),
		userTrafficDownlink: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "user_traffic_downlink"),
			"User downlink.",
			[]string{"user"}, nil,
		),

		numGoroutine: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "num_goroutine"), "", nil, nil),
		numGC:        prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "num_gc"), "", nil, nil),
		alloc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "alloc"), "", nil, nil),
		totalAlloc:   prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "total_alloc"), "", nil, nil),
		sys:          prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "sys"), "", nil, nil),
		mallocs:      prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "mallocs"), "", nil, nil),
		frees:        prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "frees"), "", nil, nil),
		liveObjects:  prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "live_objects"), "", nil, nil),
		pauseTotalNs: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "pause_total_ns"), "", nil, nil),
		uptime:       prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "uptime"), "", nil, nil),
	}
}
