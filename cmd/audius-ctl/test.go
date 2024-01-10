package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	testCmd = &cobra.Command{
		Use:   "test [command]",
		Short: "test audius-d connectivity",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	testContextCmd = &cobra.Command{
		Use:   "context",
		Short: "Test the health of different contexts",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx_config, err := conf.ReadOrCreateContextConfig()
			if err != nil {
				return logger.Error("Failed to retrieve context. ", err)
			}

			var hosts [][]string

			if len(ctx_config.CreatorNodes) > 0 {
				for _, cc := range ctx_config.CreatorNodes {
					hosts = append(hosts, []string{cc.Host, ".data.healthy"})
				}
			}
			if len(ctx_config.DiscoveryNodes) > 0 {
				for _, cc := range ctx_config.DiscoveryNodes {
					hosts = append(hosts, []string{cc.Host, ".data.discovery_provider_healthy"})
				}
			}
			if len(ctx_config.IdentityService) > 0 {
				for _, cc := range ctx_config.IdentityService {
					hosts = append(hosts, []string{cc.Host, ".healthy"})
				}
			}

			for _, host := range hosts {
				if err := checkHealth(host[0], host[1]); err != nil {
					return err
				}
			}
			return nil
		},
	}
)

type HealthCheckResponse struct {
	Data struct {
		Healthy                  bool `json:"healthy,omitempty"`
		DiscoveryProviderHealthy bool `json:"discovery_provider_healthy,omitempty"`
	} `json:"data"`
	Healthy bool `json:"healthy,omitempty"`
}

func checkHealth(host, key string) error {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	timeout := 0
	var resp *http.Response
	var err error
	for timeout < 60 {
		resp, err = httpClient.Get(host + "/health_check")
		if err != nil || resp.StatusCode == 502 {
			time.Sleep(3 * time.Second)
			timeout += 3
			fmt.Printf("Waiting for server to start... %d\n", timeout)
			continue
		}
		break
	}

	if timeout >= 60 {
		return fmt.Errorf("timed out waiting for server to start")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var healthResponse HealthCheckResponse
	if err := json.Unmarshal(body, &healthResponse); err != nil {
		fmt.Println(healthResponse)
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
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
		return fmt.Errorf("unknown key: %s", key)
	}

	fmt.Printf("%-28s [ /health_check %-32s ] %t\n", host, key, result)
	return nil
}

func init() {
	testCmd.AddCommand(testContextCmd)
}
