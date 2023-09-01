package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

var envFilePath string

func main() {
	flag.StringVar(&envFilePath, "c", "", "Path to the .env file")

	cmdName := "up"
	if len(os.Args) > 1 {
		cmdName = os.Args[1]
	}

	flag.Parse()

	switch cmdName {
	case "up":
		checkConfigFile()
		runUp()
	case "down":
		runDown()
	default:
		checkConfigFile()
		runUp()
	}
}

func checkConfigFile() {
	if envFilePath == "" {
		usr, err := user.Current()
		if err == nil {
			defaultConfig := filepath.Join(usr.HomeDir, ".audius", "audius.conf")
			if _, err := os.Stat(defaultConfig); !os.IsNotExist(err) {
				envFilePath = defaultConfig
				fmt.Println("Using default config at", defaultConfig)
			} else {
				fmt.Printf("Config not found at default location:\n\t%s\n", defaultConfig)
				fmt.Println("\nPlace your config there, or provide a valid config using the -c flag.")
				fmt.Println("\ti.e ./audius -c audius.conf")
				fmt.Println("\n# minimum required audius.conf")
				fmt.Println("creatorNodeEndpoint=")
				fmt.Println("delegateOwnerWallet=")
				fmt.Println("delegatePrivateKey=")
				fmt.Println("spOwnerWallet=")
				os.Exit(1)
			}
		} else {
			fmt.Println("Error retrieving current user:", err)
			os.Exit(1)
		}
	}
}

func runUp() {

	pullCmd := exec.Command("docker", "pull", "endliine/audius-docker-compose:linux")
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr

	if err := pullCmd.Run(); err != nil {
		fmt.Println("Error pulling image:", err)
		os.Exit(1)
	}

	volumeFlag := ""
	if envFilePath != "" {
		volumeFlag = fmt.Sprintf("-v %s:/root/audius-docker-compose/creator-node/override.env", envFilePath)
	}

	cmd := fmt.Sprintf(`docker run \
    --privileged -d \
    --name audius-creator-node \
    %s \
    -v /var/k8s/mediorum:/var/k8s/mediorum \
    -v /var/k8s/creator-node-backend:/var/k8s/creator-node-backend \
    -v /var/k8s/creator-node-db:/var/k8s/creator-node-db \
    -p 80:80 \
    -p 443:443 \
    endliine/audius-docker-compose:linux \
    && \
    docker exec audius-creator-node sh -c "while ! docker ps &> /dev/null; do echo 'docker in docker is starting up' && sleep 1; done && docker compose up -d"`, volumeFlag)

	execCommand := exec.Command("/bin/sh", "-c", cmd)
	execCommand.Stdout = os.Stdout
	execCommand.Stderr = os.Stderr

	if err := execCommand.Run(); err != nil {
		fmt.Println("Error executing command:", err)
		removeContainer()
		os.Exit(1)
	}
}

func runDown() {
	removeContainer()
}

func removeContainer() {
	fmt.Println("Removing container")
	removeCmd := exec.Command("docker", "rm", "-f", "audius-creator-node")
	if err := removeCmd.Run(); err != nil {
		fmt.Println("Error removing container:", err)
	}
}
