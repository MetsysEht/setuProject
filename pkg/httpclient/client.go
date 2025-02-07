package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/MetsysEht/setuProject/internal/config"
	"github.com/MetsysEht/setuProject/pkg/logger"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/hystrix"
	"github.com/gojek/heimdall/v7/plugins"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"golang.org/x/net/http2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InitializeClient(command string, connPoolConfig config.ConnPoolConfig, resiliencyConfig config.HystrixResiliencyConfig, retriable heimdall.Retriable, retryCount int) *hystrix.Client {

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			KeepAlive: time.Duration(connPoolConfig.KeepAliveTimeout) * time.Millisecond,
		}).DialContext,
		MaxIdleConnsPerHost:   connPoolConfig.MaxIdleConnections,
		MaxIdleConns:          connPoolConfig.MaxIdleConnections,
		IdleConnTimeout:       time.Duration(connPoolConfig.KeepAliveTimeout) * time.Millisecond,
		TLSHandshakeTimeout:   time.Duration(connPoolConfig.Timeout) * time.Millisecond,
		ExpectContinueTimeout: time.Duration(connPoolConfig.Timeout) * time.Millisecond,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: connPoolConfig.SkipCertVerification,
		},
	}

	var rt http.RoundTripper = transport

	_ = http2.ConfigureTransport(transport)

	client := &http.Client{
		Transport: &nethttp.Transport{RoundTripper: rt},
	}

	clientHystrix := hystrix.NewClient(
		hystrix.WithHTTPClient(client),
		hystrix.WithHTTPTimeout(time.Duration(connPoolConfig.Timeout)*time.Millisecond),
		hystrix.WithCommandName(command),
		hystrix.WithHystrixTimeout(time.Duration(resiliencyConfig.CircuitBreakerTimeout)*time.Millisecond),
		hystrix.WithMaxConcurrentRequests(resiliencyConfig.MaxConcurrentRequests),
		hystrix.WithErrorPercentThreshold(resiliencyConfig.ErrorPercentThreshold),
		hystrix.WithFallbackFunc(nil),
		hystrix.WithRetrier(retriable),
		hystrix.WithRetryCount(retryCount),
		hystrix.WithSleepWindow(resiliencyConfig.CircuitBreakerSleepWindow),
		hystrix.WithRequestVolumeThreshold(resiliencyConfig.RequestVolumeThreshold),
	)
	requestLogger := plugins.NewRequestLogger(nil, nil)
	clientHystrix.AddPlugin(requestLogger)
	return clientHystrix
}

func SendRequest(ctx context.Context, url string, httpMethod string, auth *config.Auth, jsonBody interface{}, headers map[string]string, doer heimdall.Doer) (string, error, int) {
	logger.L.Infow("Making http Request", httpMethod, url)
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, getIOReaderFromInterface(jsonBody))
	if err != nil {
		return "", status.Error(codes.Internal, err.Error()), http.StatusInternalServerError
	}
	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	setHeader(req, map[string]string{
		"Content-Type": "application/json",
	})
	setHeader(req, headers)
	res, err := doer.Do(req)
	if err != nil {
		logger.L.Errorw("http error", err.Error())
		return "", status.Error(codes.Internal, err.Error()), http.StatusInternalServerError
	}
	statusCode := res.StatusCode
	logger.L.Infow("http Response", statusCode)
	return getResponseBodyAsString(res), nil, statusCode
}

func getIOReaderFromInterface(data interface{}) io.Reader {
	jsonData, _ := json.Marshal(data)
	return bytes.NewReader(jsonData)
}

func setHeader(r *http.Request, headers map[string]string) {
	for key, value := range headers {
		r.Header.Set(key, value)
	}
}

func getResponseBodyAsString(resp *http.Response) string {
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	return string(bodyBytes)
}
