package health

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
)

var (
	httpClient = &http.Client{
		Timeout: time.Second * 3,
	}
)

const (
	maxRetries    = 2
	retryInterval = 1 * time.Second
)

type discoveryHealthCheckResponse struct {
	Data struct {
		Healthy                  bool `json:"healthy,omitempty"`
		DiscoveryProviderHealthy bool `json:"discovery_provider_healthy,omitempty"`
		BlockDifference          int64
		ChainHealth              struct {
			Status string `json:"status,omitempty"`
		} `json:"chain_health,omitempty"`
		FilesystemSize uint64   `json:"filesystem_size,omitempty"`
		FilesystemUsed uint64   `json:"filesystem_used,omitempty"`
		DatabaseSize   uint64   `json:"database_size,omitempty"`
		Errors         []string `json:"errors,omitempty"`
	} `json:"data,omitempty"`
	Comms struct {
		Booted time.Time `json:"booted,omitempty"`
	} `json:"comms,omitempty"`
}

type contentHealthCheckResponse struct {
	Data struct {
		Healthy          bool      `json:"healthy,omitempty"`
		MediorumPathUsed uint64    `json:"mediorumPathUsed,omitempty"`
		MediorumPathSize uint64    `json:"mediorumPathSize,omitempty"`
		DatabaseSize     uint64    `json:"databaseSize,omitempty"`
		StartedAt        time.Time `json:"startedAt,omitempty"`
	} `json:"data,omitempty"`
}

type identityHealthCheckResponse struct {
	Healthy bool `json:"healthy,omitempty"`
}

type DevnetHost struct {
	Host    string
	Healthy bool
}

type DevnetHealthSummary struct {
	Hosts []DevnetHost
}

type NodeHealthSummary struct {
	Type               conf.NodeType
	Up                 bool
	Healthy            bool
	ChainHealthy       bool
	ChainPortOpen      bool
	DatabaseSizeBytes  uint64
	DiskSpaceUsedBytes uint64
	DiskSpaceSizeBytes uint64
	BootTime           time.Time
	Errors             []string
}

func CheckDevnetHealth() DevnetHealthSummary {
	healthSummary := DevnetHealthSummary{}

	devnetDeps := map[string][]byte{
		"http://acdc-ganache.devnet.audius-d":          []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`),
		"http://eth-ganache.devnet.audius-d":           []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`),
		"http://solana-test-validator.devnet.audius-d": []byte(`{"jsonrpc":"2.0","id":1,"method":"getFirstAvailableBlock"}`),
	}
	for host, jsonData := range devnetDeps {
		hostHealth := DevnetHost{
			Host: host,
		}
		resp, err := requestWithRetries(host, jsonData)
		if err == nil {
			resp.Body.Close()
			hostHealth.Healthy = true
		}
		healthSummary.Hosts = append(healthSummary.Hosts, hostHealth)
	}
	return healthSummary
}

func CheckNodeHealth(host string, config conf.NodeConfig) (NodeHealthSummary, error) {
	switch config.Type {
	case conf.Discovery:
		return getDiscoveryNodeHealth(host, config)
	case conf.Content:
		return getContentNodeHealth(host, config)
	case conf.Identity:
		return getIdentityNodeHealth(host, config)
	default:
		return NodeHealthSummary{}, logger.Error("Unsupported node type")
	}
}

func getDiscoveryNodeHealth(host string, config conf.NodeConfig) (NodeHealthSummary, error) {
	healthSummary := NodeHealthSummary{Type: config.Type}

	resp, err := requestWithRetries(fmt.Sprintf("https://%s/health_check", host), nil)
	if err != nil {
		return healthSummary, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return healthSummary, logger.Errorf("failed to read response body: %v", err)
	}

	healthResponse := discoveryHealthCheckResponse{}
	if err := json.Unmarshal(body, &healthResponse); err != nil {
		return healthSummary, logger.Errorf("failed to parse health: %v", err)
	}
	healthSummary.Up = true
	healthSummary.Healthy = healthResponse.Data.DiscoveryProviderHealthy
	healthSummary.ChainHealthy = (strings.ToLower(healthResponse.Data.ChainHealth.Status) == "healthy")
	healthSummary.Errors = healthResponse.Data.Errors
	healthSummary.DatabaseSizeBytes = healthResponse.Data.DatabaseSize
	healthSummary.DiskSpaceUsedBytes = healthResponse.Data.FilesystemUsed
	healthSummary.DiskSpaceSizeBytes = healthResponse.Data.FilesystemSize
	healthSummary.BootTime = healthResponse.Comms.Booted

	// check if chain peering is working
	chainResp, err := requestWithRetries(fmt.Sprintf("https://%s/chain/peer", host), nil)
	if err != nil || chainResp.StatusCode == 500 {
		healthSummary.ChainPortOpen = false
	} else {
		healthSummary.ChainPortOpen = true
	}
	if chainResp != nil && chainResp.Body != nil {
		chainResp.Body.Close()
	}

	return healthSummary, nil
}

func getContentNodeHealth(host string, config conf.NodeConfig) (NodeHealthSummary, error) {
	healthSummary := NodeHealthSummary{Type: config.Type}

	resp, err := requestWithRetries(fmt.Sprintf("https://%s/health_check", host), nil)
	if err != nil {
		return healthSummary, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return healthSummary, logger.Errorf("failed to read response body: %v", err)
	}

	healthResponse := contentHealthCheckResponse{}
	if err := json.Unmarshal(body, &healthResponse); err != nil {
		return healthSummary, logger.Errorf("failed to parse health: %v", err)
	}
	healthSummary.Up = true
	healthSummary.Healthy = healthResponse.Data.Healthy
	healthSummary.DatabaseSizeBytes = healthResponse.Data.DatabaseSize
	healthSummary.DiskSpaceUsedBytes = healthResponse.Data.MediorumPathUsed
	healthSummary.DiskSpaceSizeBytes = healthResponse.Data.MediorumPathSize
	healthSummary.BootTime = healthResponse.Data.StartedAt
	return healthSummary, nil
}

func getIdentityNodeHealth(host string, config conf.NodeConfig) (NodeHealthSummary, error) {
	healthSummary := NodeHealthSummary{Type: config.Type}

	resp, err := requestWithRetries(fmt.Sprintf("https://%s/health_check", host), nil)
	if err != nil {
		return healthSummary, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return healthSummary, logger.Errorf("failed to read response body: %v", err)
	}

	healthResponse := identityHealthCheckResponse{}
	if err := json.Unmarshal(body, &healthResponse); err != nil {
		return healthSummary, logger.Errorf("failed to parse health: %v", err)
	}
	healthSummary.Up = true
	healthSummary.Healthy = healthResponse.Healthy
	return healthSummary, nil
}

func requestWithRetries(endpoint string, postData []byte) (*http.Response, error) {
	retries := 0
	var resp *http.Response
	var err error
	for retries < maxRetries {
		if postData == nil {
			resp, err = httpClient.Get(endpoint)
		} else {
			resp, err = httpClient.Post(endpoint, "application/json", bytes.NewBuffer(postData))
		}
		if err != nil || resp.StatusCode == 502 {
			time.Sleep(retryInterval)
			retries += 1
			continue
		}
		break
	}

	if retries >= maxRetries {
		return resp, logger.Errorf("Unreachable after %d retries", retries)
	}
	return resp, nil
}
