package orchestration

import (
	"fmt"
	"net/http"
	"time"
)

func awaitHealthy(containerName, host string, port uint) {
	tries := 30

	for tries > 0 {
		url := fmt.Sprintf("%s:%d/health_check", host, port)
		resp, err := http.Get(url)

		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Printf("service: %s not ready yet\n", containerName)
			time.Sleep(3 * time.Second)
			tries--
			continue
		}

		fmt.Printf("service: %s is healthy! 🎸\n", containerName)
		if resp != nil {
			resp.Body.Close()
		}
		return
	}

	fmt.Printf("%s never got healthy\n", containerName)
}
