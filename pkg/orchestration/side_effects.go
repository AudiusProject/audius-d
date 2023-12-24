package orchestration

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AudiusProject/audius-d/pkg/logger"
)

func awaitHealthy(containerName, host string, port uint) {
	tries := 30

	for tries > 0 {
		// url := fmt.Sprintf("%s:%d/health_check", host, port)
		url := fmt.Sprintf("%s/health_check", host)
		resp, err := http.Get(url)

		if err != nil || resp.StatusCode != http.StatusOK {
			logger.Infof("service: %s not ready yet\n", containerName)
			time.Sleep(3 * time.Second)
			tries--
			continue
		}

		logger.Infof("service: %s is healthy! 🎸\n", containerName)
		if resp != nil {
			resp.Body.Close()
		}
		return
	}

	logger.Infof("%s never got healthy\n", containerName)
}
