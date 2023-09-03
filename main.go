package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

//go:embed audius.conf
var confExample string
var confFilePath string

func main() {
	flag.StringVar(&confFilePath, "c", "", "Path to the .conf file")

	cmdName := "up"
	if len(os.Args) > 1 {
		cmdName = os.Args[1]
	}

	flag.Parse()

	switch cmdName {
	case "up":
		runUp(checkConfigFile())
	case "down":
		runDown()
	default:
		runUp(checkConfigFile())
	}
}

func checkConfigFile() string {
	nodeType := "discovery-provider"

	// If no config file path is provided, try to set it to the default location
	if confFilePath == "" {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Error retrieving current user:", err)
			os.Exit(1)
		}
		confFilePath = filepath.Join(usr.HomeDir, ".audius", "audius.conf")
	}

	// Check if the config file exists
	if _, err := os.Stat(confFilePath); !os.IsNotExist(err) {
		fmt.Println("Using config at", confFilePath)

		file, err := os.Open(confFilePath)
		if err != nil {
			fmt.Println("Error opening config file:", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "creatorNodeEndpoint") {
				nodeType = "creator-node"
				break
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}

	} else {
		fmt.Printf("Config not found at provided location:\n\t%s\n", confFilePath)
		fmt.Println("\nPlace your config there, or provide a valid config using the -c flag.")
		fmt.Println("\ti.e ./audius -c audius.conf")
		fmt.Printf("\n# minimum required .conf\n%s\n", confExample)
		os.Exit(1)
	}

	return nodeType
}

func runUp(nodeType string) {

	pullCmd := exec.Command("docker", "pull", "endliine/audius-docker-compose:linux")
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr

	if err := pullCmd.Run(); err != nil {
		fmt.Println("Error pulling image:", err)
		os.Exit(1)
	}

	volumeFlag := ""
	if confFilePath != "" {
		volumeFlag = fmt.Sprintf("-v %s:/root/audius-docker-compose/%s/override.env", confFilePath, nodeType)
	}

	dirPath := "/tmp/dind"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", dirPath, err)
			os.Exit(1)
		}
	}

	var cmd string

	if nodeType == "creator-node" {
		cmd = fmt.Sprintf(`docker run \
    --privileged -d \
    --name audius-creator-node \
	-v /tmp/dind:/var/lib/docker \
    %s \
    -v /var/k8s/mediorum:/var/k8s/mediorum \
    -v /var/k8s/creator-node-backend:/var/k8s/creator-node-backend \
    -v /var/k8s/creator-node-db:/var/k8s/creator-node-db \
    -p 80:80 \
    -p 443:443 \
    endliine/audius-docker-compose:linux \
    && \
    docker exec audius-creator-node sh -c "while ! docker ps &> /dev/null; do echo 'starting up' && sleep 1; done && cd %s && docker compose up -d"`, volumeFlag, nodeType)
	} else {
		cmd = fmt.Sprintf(`docker run \
	--privileged -d \
	--name audius-discovery-provider \
	-v /tmp/dind:/var/lib/docker \
	%s \
	-v /var/k8s/discovery-provider-db:/var/k8s/discovery-provider-db \
	-v /var/k8s/discovery-provider-chain:/var/k8s/discovery-provider-chain \
	-v /home/ubuntu/audius-docker-compose/discovery-provider/chain:/root/audius-docker-compose/discovery-provider/chain \
	-p 80:80 \
	-p 443:443 \
	-p 5000:5000 \
	-p 30300:30300 \
    -p 30300:30300/udp \
	endliine/audius-docker-compose:linux \
	&& \
	docker exec audius-discovery-provider sh -c "while ! docker ps &> /dev/null; do echo 'starting up' && sleep 1; done && cd %s && docker compose up -d"`, volumeFlag, nodeType)
	}

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
	removeCmd := exec.Command("docker", "rm", "-f", "audius-creator-node", "audius-discovery-provider")
	if err := removeCmd.Run(); err != nil {
		fmt.Println("Error removing container:", err)
	}
}
