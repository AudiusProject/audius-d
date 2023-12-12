package orchestration

import (
	"fmt"
	"net/http"
	"time"
)

func awaitHealthy(containerName, host string, port uint) {
	tries := 10

	for tries > 0 {
		url := fmt.Sprintf("%s:%d/health_check", host, port)
		resp, err := http.Get(url)

		if resp.StatusCode != http.StatusOK || err != nil {
			resp.Body.Close()
			fmt.Printf("service: %s not ready yet\n", containerName)
			time.Sleep(3 * time.Second)
			tries--
			continue
		}

		fmt.Printf("service: %s is healthy! 🎸\n", containerName)
		resp.Body.Close()
		return
	}

	fmt.Printf("%s never got healthy", containerName)
}
