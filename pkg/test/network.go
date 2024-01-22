package test

import (
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

	if len(ctxConfig.IdentityService) > 0 {
		for _, cc := range ctxConfig.IdentityService {
			hosts = append(hosts, []string{cc.Host, ".healthy"})
		}
	}
	if len(ctxConfig.CreatorNodes) > 0 {
		for _, cc := range ctxConfig.CreatorNodes {
			hosts = append(hosts, []string{cc.Host, ".data.healthy"})
		}
	}
	if len(ctxConfig.DiscoveryNodes) > 0 {
		for _, cc := range ctxConfig.DiscoveryNodes {
			hosts = append(hosts, []string{cc.Host, ".data.discovery_provider_healthy"})
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
		resp, err = httpClient.Get(host + "/health_check")
		if err != nil || resp.StatusCode == 502 {
			time.Sleep(3 * time.Second)
			retries += 1
			continue
		}
		break
	}

	if retries >= 3 {
		return HealthCheckResponse{}, fmt.Errorf("timed out waiting for server to start")
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
