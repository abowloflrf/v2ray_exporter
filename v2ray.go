package main

import (
	"context"
	"time"

	v2stats "github.com/v2fly/v2ray-core/v4/app/stats/command"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	v2stats.StatsServiceClient
}

func NewClient(addr string) (*Client, error) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := v2stats.NewStatsServiceClient(conn)
	return &Client{
		conn:               conn,
		StatsServiceClient: client,
	}, nil
}

func (c *Client) Stats(pt string) ([]*v2stats.Stat, error) {
	start := time.Now()
	resp, err := c.QueryStats(context.Background(), &v2stats.QueryStatsRequest{
		Pattern: "",
		Reset_:  false,
	})
	if err != nil {
		return nil, err
	}
	logger.Debugln("QueryStats from v2ray", "duration", time.Since(start))
	return resp.Stat, nil
}

// SysStats
// NumGoroutine:17
// NumGC:12
// Alloc:3195136
// TotalAlloc:17653480
// Sys:71500024
// Mallocs:369427
// Frees:347892
// LiveObjects:21535
// PauseTotalNs:400819
// Uptime:1057
func (c *Client) SysStats() (*v2stats.SysStatsResponse, error) {
	start := time.Now()
	resp, err := c.GetSysStats(context.Background(), &v2stats.SysStatsRequest{})
	if err != nil {
		return nil, err
	}
	logger.Debugln("GetSysStats from v2ray", "duration", time.Since(start))
	return resp, nil
}

func (c *Client) Close() {
	logger.Warn("v2ray grpc connection closing")
	err := c.conn.Close()
	if err != nil {
		logger.Warnln("close v2ray connection", err.Error())
	}
	logger.Warn("v2ray grpc connection closed")
}
