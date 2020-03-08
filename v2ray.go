package main

import (
	"context"
	"time"

	"google.golang.org/grpc"
	v2Stats "v2ray.com/core/app/stats/command"
)

type Client struct {
	conn   *grpc.ClientConn
	client v2Stats.StatsServiceClient
}

func NewClient(addr string) (*Client, error) {
	c := new(Client)
	var err error
	c.conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c.client = v2Stats.NewStatsServiceClient(c.conn)
	return c, nil
}

func (c *Client) QueryStats(pt string) ([]*v2Stats.Stat, error) {
	start := time.Now()
	resp, err := c.client.QueryStats(context.Background(), &v2Stats.QueryStatsRequest{
		Pattern: "",
		Reset_:  false,
	})
	if err != nil {
		return nil, err
	}
	sugar.Debugw("QueryStats from v2ray", "duration", time.Since(start))
	return resp.Stat, nil
}

// GetSysStats
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
func (c *Client) GetSysStats() (*v2Stats.SysStatsResponse, error) {
	start := time.Now()
	resp, err := c.client.GetSysStats(context.Background(), &v2Stats.SysStatsRequest{})
	if err != nil {
		return nil, err
	}
	sugar.Infow("GetSysStats from v2ray", "duration", time.Since(start))
	return resp, nil
}

func (c *Client) Close() {
	sugar.Warn("v2ray grpc connection closing")
	err := c.conn.Close()
	if err != nil {
		sugar.Warnw("close v2ray connection", "error", err.Error())
	}
	sugar.Warn("v2ray grpc connection closed")
}
