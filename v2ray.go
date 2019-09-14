package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	v2Stats "v2ray.com/core/app/stats/command"
)

type Client struct {
	conn   *grpc.ClientConn
	client v2Stats.StatsServiceClient
}

func NewClient(addr string) *Client {
	c := new(Client)
	var err error
	c.conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dial gRPC server error: " + err.Error())
	}
	c.client = v2Stats.NewStatsServiceClient(c.conn)
	return c
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
	log.Println("Query from v2ray, duration:", time.Now().Sub(start))
	return resp.Stat, nil
}

func (c *Client) Close() {
	log.Println("V2ray grpc connection closing")
	err := c.conn.Close()
	if err != nil {
		log.Println("Error when close v2ray connection: " + err.Error())
	}
	log.Println("V2ray grpc connection closed")
}
