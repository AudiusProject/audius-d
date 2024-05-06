package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AudiusProject/audius-d/pkg/conf"
)

type HealthCheckResponse struct {
	Host string
	Data struct {
		Healthy                  bool `json:"healthy,omitempty"`
		DiscoveryProviderHealthy bool `json:"discovery_provider_healthy,omitempty"`
	} `json:"data"`
	Healthy bool `json:"healthy,omitempty"`
	Key     string
	Result  bool
	Error   error
}

func CheckNodeHealth(ctxConfig *conf.ContextConfig) ([]HealthCheckResponse, error) {
	var hosts [][]string
	var responses []HealthCheckResponse

	if ctxConfig.Network.DeployOn == "devnet" {
		devnetDeps := map[string][]byte{
			"http://acdc-ganache.devnet.audius-d":          []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`),
			"http://eth-ganache.devnet.audius-d":           []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`),
			"http://solana-test-validator.devnet.audius-d": []byte(`{"jsonrpc":"2.0","id":1,"method":"getFirstAvailableBlock"}`),
		}
		for host, jsonData := range devnetDeps {
			healthResponse, err := checkRPCHealth(host, jsonData)
			if err != nil {
				healthResponse = HealthCheckResponse{
					Host:  host,
					Error: err,
				}
			}
			responses = append(responses, healthResponse)
		}
	}

	for host, config := range ctxConfig.Nodes {
		switch config.Type {
		case conf.Identity:
			hosts = append(hosts, []string{host, ".healthy"})
		case conf.Content:
			hosts = append(hosts, []string{host, ".data.healthy"})
		case conf.Discovery:
			hosts = append(hosts, []string{host, ".data.discovery_provider_healthy"})
		}
	}

	for _, host := range hosts {
		healthResponse, err := checkHealth(host[0], host[1])
		if err != nil {
			healthResponse = HealthCheckResponse{
				Host:  host[0],
				Key:   host[1],
				Error: err,
			}
		}
		responses = append(responses, healthResponse)
	}

	return responses, nil
}

func checkHealth(host, key string) (HealthCheckResponse, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 3,
	}

	retries := 0
	var resp *http.Response
	var err error
	for retries < 3 {
		resp, err = httpClient.Get(fmt.Sprintf("https://%s/health_check", host))
		if err != nil || resp.StatusCode == 502 {
			time.Sleep(3 * time.Second)
			retries += 1
			continue
		}
		break
	}

	if retries >= 3 {
		return HealthCheckResponse{}, fmt.Errorf("timed out after %d retries", retries)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return HealthCheckResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var healthResponse HealthCheckResponse
	if err := json.Unmarshal(body, &healthResponse); err != nil {
		return HealthCheckResponse{}, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	var result bool
	switch key {
	case ".data.healthy":
		result = healthResponse.Data.Healthy
	case ".data.discovery_provider_healthy":
		result = healthResponse.Data.DiscoveryProviderHealthy
	case ".healthy":
		result = healthResponse.Healthy
	default:
		return HealthCheckResponse{}, fmt.Errorf("unknown key: %s", key)
	}

	return HealthCheckResponse{
		Host:   host,
		Key:    key,
		Result: result,
	}, nil
}

func checkRPCHealth(host string, jsonData []byte) (HealthCheckResponse, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 3,
	}

	retries := 0
	for retries < 3 {
		resp, err := httpClient.Post(host, "application/json", bytes.NewBuffer(jsonData))
		if err != nil || resp.StatusCode == 502 {
			time.Sleep(3 * time.Second)
			retries += 1
			continue
		}
		resp.Body.Close()
		break
	}

	if retries >= 3 {
		return HealthCheckResponse{Host: host, Error: fmt.Errorf("timed out after %d retries", retries)}, nil
	}

	return HealthCheckResponse{Host: host, Result: true}, nil
}
