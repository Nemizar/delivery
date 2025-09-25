package geo

import (
	"context"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/generated/clients/geosrv/geopb"
	"delivery/internal/pkg/errs"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ ports.GeoClient = &client{}

type client struct {
	conn        *grpc.ClientConn
	pbGeoClient geopb.GeoClient
	timeout     time.Duration
}

func NewClient(host string) (ports.GeoClient, error) {
	if host == "" {
		return nil, errs.NewValueIsRequiredError("host")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	pbGeoClient := geopb.NewGeoClient(conn)

	return &client{
		conn:        conn,
		pbGeoClient: pbGeoClient,
		timeout:     5 * time.Second,
	}, nil
}

func (c *client) GetLocation(ctx context.Context, street string) (kernel.Location, error) {
	req := &geopb.GetGeolocationRequest{
		Street: street,
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.pbGeoClient.GetGeolocation(ctx, req)
	if err != nil {
		return kernel.Location{}, err
	}

	return kernel.NewLocation(int(resp.Location.GetX()), int(resp.Location.GetY()))
}

func (c *client) Close() error {
	return c.conn.Close()
}
