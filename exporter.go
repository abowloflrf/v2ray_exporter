package main

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
	v2Stats "v2ray.com/core/app/stats/command"
)

const (
	namespace = "v2ray"
)

type Exporter struct {
	inboundTagTrafficUplink   *prometheus.Desc
	inboundTagTrafficDownlink *prometheus.Desc
	userTrafficUplink         *prometheus.Desc
	userTrafficDownlink       *prometheus.Desc
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
	stats, err := v2c.QueryStats("")
	if err != nil {
		log.Println("Failed to collect stats from v2ray:", err)
	}
	if err := e.parseSystemStats(ch, stats); err != nil {
		log.Println("Could not parse stats from v2ray:", err)
	}
}

func (e *Exporter) parseStatsItem(s *v2Stats.Stat) (*SingleF64Stat, error) {
	//"user>>>i@ruofeng.me>>>traffic>>>downlink"
	//"inbound>>>api>>>traffic>>>uplink"
	//"inbound>>>vmess-ws-in>>>traffic>>>downlink"
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

func (e *Exporter) parseSystemStats(ch chan<- prometheus.Metric, stats []*v2Stats.Stat) error {
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
				log.Println(err)
				break
			}
			if single.Name == m {
				ch <- prometheus.MustNewConstMetric(d, single.Type, single.Value, single.Tag)
				log.Println("Collected!", single)
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
	}
}
