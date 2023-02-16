package geek

import (
	"context"
	"fmt"
	pb "geek-cache/geek/pb"
	registry "geek-cache/geek/registry"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	name string // name of remote server, e.g. ip:port
}

// NewClient creates a new client
func NewClient(name string) (*Client, error) {
	return &Client{name: name}, nil
}

// Get send the url for getting specific group and key,
// and return the result
func (c *Client) Get(group, key string) ([]byte, error) {
	cli, err := clientv3.New(registry.DefaultEtcdConfig)

	if err != nil {
		return nil, err
	}
	defer cli.Close()

	conn, err := registry.EtcdDial(cli, c.name)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	grpcCLient := pb.NewGroupCacheClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := grpcCLient.Get(ctx, &pb.Request{
		Group: group,
		Key:   key,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get %s-%s from peer %s", group, key, c.name)
	}
	return resp.GetValue(), nil
}

// resure implemented
var _ PeerGetter = (*Client)(nil)