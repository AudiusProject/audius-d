package orchestration

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
)

func awaitHealthy(config *conf.ContextConfig) {
	for cname, cc := range config.CreatorNodes {
		awaitChan := make(chan string)
		go awaitService(cname, cc.Host, awaitChan)
		for log := range awaitChan {
			logger.Info(log)
		}
	}
	for cname, dc := range config.DiscoveryNodes {
		awaitChan := make(chan string)
		go awaitService(cname, dc.Host, awaitChan)
		for log := range awaitChan {
			logger.Info(log)
		}
	}
	for cname, id := range config.IdentityService {
		awaitChan := make(chan string)
		go awaitService(cname, id.Host, awaitChan)
		for log := range awaitChan {
			logger.Info(log)
		}
	}
}

func awaitService(containerName, host string, awaitChan chan string) {
	defer close(awaitChan)
	tries := 30

	for tries > 0 {

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		url := fmt.Sprintf("%s/health_check", host)
		resp, err := client.Get(url)

		if err != nil || resp.StatusCode != http.StatusOK {
			awaitChan <- fmt.Sprintf("service: %s not ready yet\n", url)
			time.Sleep(3 * time.Second)
			tries--
			continue
		}

		awaitChan <- fmt.Sprintf("service: %s is healthy! 🎸\n", containerName)
		if resp != nil {
			resp.Body.Close()
		}
		return
	}

	awaitChan <- fmt.Sprintf("%s never got healthy\n", containerName)
}
