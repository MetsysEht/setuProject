package setuGateway

import (
	"context"
	"encoding/json"

	"github.com/MetsysEht/setuProject/internal/config"
	"github.com/MetsysEht/setuProject/pkg/httpclient"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/hystrix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gateway struct {
	client *hystrix.Client
	cfg    config.SetuGatewayService
}

func NewGateway(cfg config.SetuGatewayService) ISetuGateway {
	client := httpclient.InitializeClient(
		"setuService",
		cfg.ConnPoolConfig,
		cfg.HystrixResiliencyConfig,
		heimdall.NewNoRetrier(),
		3,
	)
	g := gateway{client: client, cfg: cfg}
	if cfg.Mock == true {
		return &g
	}
	return &g
}

func (g *gateway) VerifyPan(ctx context.Context, request *PANRequest) (*PANResponse, error) {
	url := g.cfg.BaseUrl + g.cfg.ValidatePAN.Path

	resp, err, statusCode := httpclient.SendRequest(ctx, url, g.cfg.ValidatePAN.Method, nil, request, g.cfg.ValidatePAN.Headers, g.client)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, status.Error(codes.Internal, "API Failure")
	}
	response := PANResponse{}
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &response, nil
}
